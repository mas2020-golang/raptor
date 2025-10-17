<#
.SYNOPSIS
  Install the latest (or a specific) Raptor CLI on Windows.

.DESCRIPTION
  - Detects CPU architecture (amd64/arm64).
  - Fetches GitHub release assets from mas2020-golang/raptor.
  - Downloads the Windows ZIP (as published by GoReleaser).
  - Extracts raptor.exe and installs to $env:LOCALAPPDATA\Programs\raptor by default.
  - Adds that folder to the user's PATH if needed.

.PARAMETER Version
  Optional tag like 'v0.3.1'. If omitted, uses the latest release.

.PARAMETER InstallDir
  Optional install directory. Default: $env:LOCALAPPDATA\Programs\raptor

.PARAMETER Force
  Overwrite any existing raptor.exe in InstallDir.

.PARAMETER Quiet
  Reduce output (still shows errors).

.PARAMETER PrintTinyLoader
  Prints a tiny README-friendly one-liner that downloads and runs this script.

.PARAMETER Proxy
  Optional explicit proxy URL, e.g. http://proxy.mycorp.tld:8080

.PARAMETER ProxyCredential
  Optional PSCredential to authenticate at the proxy.
#>

[CmdletBinding()]
param(
  [string]$Version = "",
  [string]$InstallDir = "$env:LOCALAPPDATA\Programs\raptor",
  [switch]$Force,
  [switch]$Quiet,
  [switch]$PrintTinyLoader,
  [string]$Proxy = "",
  [pscredential]$ProxyCredential
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$ProgressPreference = 'SilentlyContinue'

# Ensure TLS 1.2+ (GitHub requires it)
try {
  [Net.ServicePointManager]::SecurityProtocol = `
    [Net.SecurityProtocolType]::Tls12 -bor `
    [Net.SecurityProtocolType]::Tls13 -bor `
    [Net.SecurityProtocolType]::Tls11
} catch { }

# ---------------------------
# Utilities
# ---------------------------
function Write-Info($msg) { if (-not $Quiet) { Write-Host "[INFO] $msg" -ForegroundColor Cyan } }
function Write-Warn($msg) { if (-not $Quiet) { Write-Host "[WARN] $msg" -ForegroundColor Yellow } }
function Write-Err($msg)  { Write-Host "[ERR ] $msg" -ForegroundColor Red }

function Get-Arch {
  $arch = $env:PROCESSOR_ARCHITECTURE
  if (-not $arch -and $env:PROCESSOR_ARCHITEW6432) { $arch = $env:PROCESSOR_ARCHITEW6432 }
  switch (($arch | ForEach-Object { "$_".ToLower() })) {
    "amd64" { return "amd64" }
    "arm64" { return "arm64" }
    default { return "amd64" }
  }
}

function Map-NameArch($arch) { if ($arch -eq "amd64") { "x86_64" } else { $arch } }

function Ensure-Dir($dir) {
  if (-not (Test-Path -LiteralPath $dir)) { New-Item -ItemType Directory -Path $dir | Out-Null }
}

function Ensure-UserPathContains($dir) {
  $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
  $segs = @(); if ($userPath) { $segs = $userPath -split ";" }
  if ($segs -notcontains $dir) {
    $new = if ([string]::IsNullOrWhiteSpace($userPath)) { $dir } else { "$userPath;$dir" }
    [Environment]::SetEnvironmentVariable("Path", $new, "User")
    Write-Info "Added to user PATH: $dir"
    Write-Warn "Open a new terminal (or sign out/in) for PATH changes to take effect."
  }
  if (-not ($env:Path -split ";" | Where-Object { $_ -eq $dir })) { $env:Path = "$env:Path;$dir" }
}

# ----- Proxy helpers -----
function Resolve-ProxyAuto {
  # WinINET (Internet Options)
  try {
    $reg = Get-ItemProperty 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings' -ErrorAction Stop
    if ($reg.ProxyEnable -eq 1 -and $reg.ProxyServer) {
      $p = $reg.ProxyServer
      if ($p -and $p -notmatch '^\w+://') { $p = "http://$p" }
      return $p
    }
  } catch { }

  # Env fallbacks
  if ($env:HTTPS_PROXY) { return $env:HTTPS_PROXY }
  if ($env:HTTP_PROXY)  { return $env:HTTP_PROXY }
  return $null
}

function Should-BypassProxyFor([string]$Url) {
  try { $host = ([uri]$Url).Host } catch { return $false }
  $skip = $env:NO_PROXY; if (-not $skip) { $skip = $env:no_proxy }
  if (-not $skip) { return $false }
  foreach ($token in ($skip -split ',' | ForEach-Object { $_.Trim() })) {
    if ($token -and $host -like "*$token") { return $true }
  }
  return $false
}

# ----- Download (parameter-first) -----
function Download-File($url, $outPath) {
  Write-Info "Downloading: $url"

  $headers = @{ 'User-Agent' = 'raptor-installer' }
  $common = @{
    Uri                = $url
    OutFile            = $outPath
    UseBasicParsing    = $true
    Headers            = $headers
    MaximumRedirection = 8
    ErrorAction        = 'Stop'
  }

  # 0) If NO_PROXY says bypass, do plain
  if (Should-BypassProxyFor $url) {
    try { Invoke-WebRequest @common; return } catch { throw "Download failed: $url`n$($_.Exception.Message)" }
  }

  # 1) If user provided -Proxy / -ProxyCredential, prefer them
  if ($Proxy) {
    try {
      if ($ProxyCredential) {
        Invoke-WebRequest @common -Proxy $Proxy -ProxyCredential $ProxyCredential
      } else {
        Invoke-WebRequest @common -Proxy $Proxy
      }
      return
    } catch {
      $msg = $_.Exception.Message
      throw "Download failed using explicit proxy: $Proxy`n$msg"
    }
  }

  # 2) Try plain first (lets PAC/auto-proxy do its thing)
  try {
    Invoke-WebRequest @common
    return
  } catch {
    $raw1 = $_.Exception.Message
    # 3) Fallback: try auto-resolved proxy (WinINET/env)
    $autoProxy = Resolve-ProxyAuto
    if ($autoProxy) {
      Write-Warn "Retrying with detected proxy: $autoProxy"
      try {
        Invoke-WebRequest @common -Proxy $autoProxy -ProxyCredential ([pscredential]::Empty)
        return
      } catch {
        $raw2 = $_.Exception.Message
        $status2 = $null; try { $status2 = $_.Exception.Response.StatusCode.value__ } catch { }
        $needsAuth = ($status2 -eq 407) -or ($raw2 -match '(?i)407|Proxy Authentication')
        if ($needsAuth) {
          Write-Warn "Proxy requires authentication. Prompting for credentials…"
          $cred = $null
          try { $cred = Get-Credential -Message "Enter your proxy (DOMAIN\\user) credentials" } catch { }
          if ($cred) {
            try {
              Invoke-WebRequest @common -Proxy $autoProxy -ProxyCredential $cred
              return
            } catch {
              throw "Download failed after proxy auth retry: $url`n$($_.Exception.Message)"
            }
          } else {
            throw "Proxy authentication required but no credentials were provided."
          }
        }
        throw "Download failed with detected proxy ($autoProxy): $url`n$raw2"
      }
    }
    throw "Download failed: $url`n$raw1"
  }
}

function Extract-Zip($zipPath, $destDir) {
  Write-Info "Extracting: $zipPath"
  Expand-Archive -Path $zipPath -DestinationPath $destDir -Force
}

function Find-RaptorExe($root) {
  Get-ChildItem -Path $root -Recurse -File -Filter "raptor.exe" | Select-Object -First 1
}

function Print-TinyLoader {
  $api = "https://api.github.com/repos/mas2020-golang/raptor/contents/install.ps1?ref=main"
  $loader = @"
# Windows (PowerShell) — install latest Raptor (no admin needed)
`$u='$api'
`$p="$env:TEMP\raptor-install-$([guid]::NewGuid().ToString('N')).ps1"
iwr `$u -Headers @{ 'User-Agent'='raptor-installer'; 'Accept'='application/vnd.github.v3.raw' } -UseBasicParsing -OutFile `$p
Unblock-File `$p
powershell.exe -NoProfile -ExecutionPolicy Bypass -File `$p
"@
  $loader
}

if ($PrintTinyLoader) { Print-TinyLoader; return }

# ---------------------------
# Install flow
# ---------------------------
$repo = "mas2020-golang/raptor"
$apiBase = "https://api.github.com/repos/$repo/releases"
$headers = @{ "User-Agent" = "raptor-installer" }

$target = if ([string]::IsNullOrWhiteSpace($Version)) { "$apiBase/latest" } else { "$apiBase/tags/$Version" }

Write-Info "Target repo: $repo"
Write-Info ("Version: " + ($(if ($Version) { $Version } else { "latest" })))

try {
  $release = Invoke-RestMethod -Uri $target -Headers $headers
} catch {
  if ($Version) { throw "Unable to fetch release '$Version' from GitHub. $($_.Exception.Message)" }
  throw "Unable to fetch the latest release from GitHub. $($_.Exception.Message)"
}

if (-not $release.assets -or $release.assets.Count -eq 0) { throw "No assets found in release $($release.tag_name)" }

$arch = Get-Arch
$nameArch = Map-NameArch $arch
$expected = "raptor_Windows_${nameArch}.zip"

$asset = $release.assets | Where-Object { $_.name -eq $expected } | Select-Object -First 1
if (-not $asset) { $asset = $release.assets | Where-Object { $_.name -match ("(?i)^raptor[_-]windows[_-](" + [regex]::Escape($nameArch) + ")\.zip$") } | Select-Object -First 1 }
if (-not $asset) { $asset = $release.assets | Where-Object { $_.name -match "(?i)windows.*\.zip$" } | Select-Object -First 1 }
if (-not $asset) { throw "Could not find a Windows ZIP asset in release $($release.tag_name) for arch '$arch'. Expected like: $expected" }

Write-Info "Selected asset: $($asset.name)"

$tmpRoot    = Join-Path $env:TEMP ("raptor-install-" + [Guid]::NewGuid().ToString("N"))
$zipPath    = Join-Path $tmpRoot $asset.name
$extractDir = Join-Path $tmpRoot "extract"
Ensure-Dir $tmpRoot; Ensure-Dir $extractDir

Download-File $asset.browser_download_url $zipPath
Extract-Zip $zipPath $extractDir

$exe = Find-RaptorExe $extractDir
if (-not $exe) { throw "raptor.exe not found after extracting $($asset.name)" }

Ensure-Dir $InstallDir
$dest = Join-Path $InstallDir "raptor.exe"

if ((Test-Path -LiteralPath $dest) -and (-not $Force)) {
  Write-Warn "raptor.exe already exists at: $dest"
  Write-Warn "Re-run with -Force to overwrite, or specify -InstallDir to install elsewhere."
} else {
  Copy-Item -Path $exe.FullName -Destination $dest -Force
  Write-Info "Installed: $dest"
}

Ensure-UserPathContains $InstallDir

if (Test-Path -LiteralPath $dest) {
  if (-not $Quiet) {
    Write-Host ""
    Write-Host "Raptor was installed successfully." -ForegroundColor Green
    Write-Host "Open a new terminal and run:" -ForegroundColor Green
    Write-Host "  raptor version" -ForegroundColor Green
  }
} else {
  throw "Installation failed: $dest not found."
}