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

.EXAMPLES
  powershell -ExecutionPolicy Bypass -File .\install.ps1
  powershell -ExecutionPolicy Bypass -File .\install.ps1 -Version v0.3.1
  powershell -ExecutionPolicy Bypass -File .\install.ps1 -InstallDir "C:\Tools\raptor" -Force
  powershell -ExecutionPolicy Bypass -File .\install.ps1 -PrintTinyLoader
#>

[CmdletBinding()]
param(
  [string]$Version = "",
  [string]$InstallDir = "$env:LOCALAPPDATA\Programs\raptor",
  [switch]$Force,
  [switch]$Quiet,
  [switch]$PrintTinyLoader
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# ---------------------------
# Utilities
# ---------------------------
function Write-Info($msg) { if (-not $Quiet) { Write-Host "[INFO] $msg" -ForegroundColor Cyan } }
function Write-Warn($msg) { if (-not $Quiet) { Write-Host "[WARN] $msg" -ForegroundColor Yellow } }
function Write-Err($msg)  { Write-Host "[ERR ] $msg" -ForegroundColor Red }

function Get-Arch {
  # Normalize to our artifact naming
  # amd64 => x86_64 in asset name; arm64 => arm64
  $arch = $env:PROCESSOR_ARCHITECTURE
  if (-not $arch -and $env:PROCESSOR_ARCHITEW6432) { $arch = $env:PROCESSOR_ARCHITEW6432 }
  switch (($arch | ForEach-Object { "$_".ToLower() })) {
    "amd64" { return "amd64" }
    "arm64" { return "arm64" }
    default { return "amd64" } # safe default
  }
}

function Map-NameArch($arch) {
  if ($arch -eq "amd64") { return "x86_64" }
  return $arch # arm64 stays arm64
}

function Ensure-Dir($dir) {
  if (-not (Test-Path -LiteralPath $dir)) {
    New-Item -ItemType Directory -Path $dir | Out-Null
  }
}

function Ensure-UserPathContains($dir) {
  $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
  $segs = @()
  if ($userPath) { $segs = $userPath -split ";" }
  if ($segs -notcontains $dir) {
    $new = if ([string]::IsNullOrWhiteSpace($userPath)) { $dir } else { "$userPath;$dir" }
    [Environment]::SetEnvironmentVariable("Path", $new, "User")
    Write-Info "Added to user PATH: $dir"
    Write-Warn "Open a new terminal (or sign out/in) for PATH changes to take effect."
  }
  # also add to current process
  if (-not ($env:Path -split ";" | Where-Object { $_ -eq $dir })) {
    $env:Path = "$env:Path;$dir"
  }
}

function Download-File($url, $outPath) {
  Write-Info "Downloading: $url"
  try {
    Invoke-WebRequest -UseBasicParsing -Uri $url -OutFile $outPath
  } catch {
    throw "Download failed: $url`n$($_.Exception.Message)"
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
  # Use the GitHub API with Accept: raw so we always get the script body (avoids HTML/JS pages).
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

if ($PrintTinyLoader) {
  Print-TinyLoader
  return
}

# ---------------------------
# Install flow
# ---------------------------
$repo = "mas2020-golang/raptor"
$apiBase = "https://api.github.com/repos/$repo/releases"
$headers = @{ "User-Agent" = "raptor-installer" }

# --- PS5-safe target selection (no "if-as-expression") ---
if ([string]::IsNullOrWhiteSpace($Version)) {
  $target = "$apiBase/latest"
} else {
  $target = "$apiBase/tags/$Version"
}

Write-Info "Target repo: $repo"

# PS5-safe version label
$verLabel = $Version
if ([string]::IsNullOrWhiteSpace($verLabel)) { $verLabel = "latest" }
Write-Info ("Version: " + $verLabel)

# Query release JSON
try {
  $release = Invoke-RestMethod -Uri $target -Headers $headers
} catch {
  if (-not [string]::IsNullOrWhiteSpace($Version)) {
    throw "Unable to fetch release '$Version' from GitHub. $($_.Exception.Message)"
  } else {
    throw "Unable to fetch the latest release from GitHub. $($_.Exception.Message)"
  }
}

if (-not $release.assets -or $release.assets.Count -eq 0) {
  throw "No assets found in release $($release.tag_name)"
}

# Choose asset based on our naming: raptor_Windows_x86_64.zip | raptor_Windows_arm64.zip
$arch = Get-Arch
$nameArch = Map-NameArch $arch
$expected = "raptor_Windows_${nameArch}.zip"

$asset = $release.assets | Where-Object { $_.name -eq $expected } | Select-Object -First 1
if (-not $asset) {
  # Be a little forgiving if case varies
  $asset = $release.assets | Where-Object { $_.name -match ("(?i)^raptor[_-]windows[_-](" + [regex]::Escape($nameArch) + ")\.zip$") } | Select-Object -First 1
}
if (-not $asset) {
  # Final fallback: any windows zip (useful if naming ever changes)
  $asset = $release.assets | Where-Object { $_.name -match "(?i)windows.*\.zip$" } | Select-Object -First 1
}
if (-not $asset) {
  throw "Could not find a Windows ZIP asset in release $($release.tag_name) for arch '$arch'. Expected like: $expected"
}

Write-Info "Selected asset: $($asset.name)"

# Temp workspace
$tmpRoot = Join-Path $env:TEMP ("raptor-install-" + [Guid]::NewGuid().ToString("N"))
Ensure-Dir $tmpRoot
$zipPath = Join-Path $tmpRoot $asset.name
$extractDir = Join-Path $tmpRoot "extract"
Ensure-Dir $extractDir

# Download and extract
Download-File $asset.browser_download_url $zipPath
Extract-Zip $zipPath $extractDir

# Locate raptor.exe
$exe = Find-RaptorExe $extractDir
if (-not $exe) {
  throw "raptor.exe not found after extracting $($asset.name)"
}

# Install
Ensure-Dir $InstallDir
$dest = Join-Path $InstallDir "raptor.exe"

if ((Test-Path -LiteralPath $dest) -and (-not $Force)) {
  Write-Warn "raptor.exe already exists at: $dest"
  Write-Warn "Re-run with -Force to overwrite, or specify -InstallDir to install elsewhere."
} else {
  Copy-Item -Path $exe.FullName -Destination $dest -Force
  Write-Info "Installed: $dest"
}

# PATH
Ensure-UserPathContains $InstallDir

# Final check
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