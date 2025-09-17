# Raptor <!-- omit in toc -->

Raptor is a cross-platform CLI application for managing your secrets in a fast and smart way.  
It lets you **encrypt and decrypt files or folders** and manage **encrypted boxes** (JSON-based containers) where you can safely store credentials, tokens, and other sensitive data.  

Raptor is designed to be simple, portable, and secure — working on macOS, Linux, and Windows.

---
- [Installing Raptor](#installing-raptor)
  - [From Source (for Linux and Mac users)](#from-source-for-linux-and-mac-users)
  - [macOS / Linux](#macos--linux)
  - [Optional dependencies (Linux clipboard)](#optional-dependencies-linux-clipboard)
  - [Windows (PowerShell)](#windows-powershell)
  - [Uninstall on macOS/Linux](#uninstall-on-macoslinux)
  - [Uninstall on Windows](#uninstall-on-windows)
  - [Troubleshooting](#troubleshooting)
- [Command Reference](#command-reference)
- [Typical Workflow](#typical-workflow)
- [Usage Examples](#usage-examples)
  - [Encrypt a File](#encrypt-a-file)
  - [Decrypt a File](#decrypt-a-file)
  - [Create a Box](#create-a-box)
  - [Add a Secret to a Box](#add-a-secret-to-a-box)
  - [Generate a Random Password](#generate-a-random-password)
  - [List Secrets in a Box](#list-secrets-in-a-box)
  - [Get a Secret and Copy to Clipboard](#get-a-secret-and-copy-to-clipboard)
  - [Edit a Secret](#edit-a-secret)
- [Environment Variables](#environment-variables)
- [How It Works](#how-it-works)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Installing Raptor

### From Source (for Linux and Mac users)
You need Go 1.21 or later.

```bash
git clone https://github.com/your/repo.git
cd raptor
make build
```
This will produce a binary `raptor` in the project folder.

### macOS / Linux
Use the installer script to download the latest release for your OS/CPU and place the binary on your PATH.
```shell
# macOS or Linux
curl -fsSL https://raw.githubusercontent.com/mas2020-golang/raptor/main/install.sh | bash
```
This script:

- Detects your OS/architecture
- Downloads the latest release
- Installs raptor into a standard location (e.g., /usr/local/bin or ~/.local/bin)
- Makes it executable

Verify
```shell
raptor version
raptor help
```

### Optional dependencies (Linux clipboard)
On Linux, clipboard features (copying secrets directly) require one of these utilities:
- `xclip`
- `xsel`

Install them via your package manager if you want clipboard integration.

### Windows (PowerShell)
```shell
iwr -UseBasicParsing https://raw.githubusercontent.com/mas2020-golang/raptor/main/install.ps1 | iex
```

### Uninstall on macOS/Linux
Remove the binary (adjust the path if needed):
```shell
sudo rm -f /usr/local/bin/raptor
# or, if installed to ~/.local/bin
rm -f ~/.local/bin/raptor
```
If you have any box installed, remove it:
```shell
rm -rf $HOME/.cryptex/boxes
```

### Uninstall on Windows
Remove the binary (adjust the path if needed):

- delete `%LOCALAPPDATA%\Programs\raptor\raptor.exe`
- optionally remove `%LOCALAPPDATA%\Programs\raptor` from your User PATH.

### Troubleshooting

- `raptor`: command not found
Ensure the install location is on your PATH. Common paths:
   - macOS (Intel/Apple Silicon): `/usr/local/bin` or `/opt/homebrew/bin`
   - Linux: `/usr/local/bin` or `~/.local/bin`
- Clipboard flags don’t copy on Linux
Install `xclip` or `xsel` (see “Optional dependencies”). 

---

## Command Reference

| Command | Description |
|---------|-------------|
| `raptor encrypt --in FILE --out FILE` | Encrypt a file |
| `raptor decrypt --in FILE --out FILE` | Decrypt a file |
| `raptor create box --name NAME` | Create a new box |
| `raptor create secret --box NAME --name KEY --value VAL` | Add a secret to a box |
| `raptor create password [--length N] [--symbols]` | Generate a random password |
| `raptor list box` | List all existing boxes |
| `raptor list secret --box NAME` | List secrets in a box |
| `raptor get secret --box NAME --name KEY [--clip]` | Retrieve a secret (optionally copy to clipboard) |
| `raptor edit secret --box NAME --name KEY` | Edit a secret in the default editor |
| `raptor print box --box NAME` | Print all secrets in a box |
| `raptor open --box NAME` | Open a box and keep it active until timeout |
| `raptor version` | Show Raptor version info |

Run `raptor help <command>` for full details on options.

---

## Typical Workflow

Here’s a quick journey through Raptor’s main features:

1. **Encrypt a file you want to protect**  
   ```bash
   raptor encrypt --in secrets.env --out secrets.env.enc
   ```

2. **Create a box to organize secrets**  
   ```bash
   raptor create box --name my-box
   ```

3. **Add secrets into the box**  
   ```bash
   raptor create secret --box my-box --name DB_PASSWORD --value "SuperSecret123"
   raptor create secret --box my-box --name API_KEY --value "sk_live_abc123"
   ```

4. **List the secrets in your box**  
   ```bash
   raptor list secret --box my-box
   ```

5. **Retrieve a secret safely**  
   ```bash
   raptor get secret --box my-box --name API_KEY --clip
   ```

6. **Edit or update a secret when it changes**  
   ```bash
   raptor edit secret --box my-box --name DB_PASSWORD
   ```

7. **Print or open a box when you need to work interactively**  
   ```bash
   raptor open --box my-box
   ```

This flow covers the most common tasks: protecting files, creating secure containers, and handling credentials.

---

## Usage Examples

### Encrypt a File
```bash
raptor encrypt --in secrets.env --out secrets.env.enc
```

### Decrypt a File
```bash
raptor decrypt --in secrets.env.enc --out secrets.env
```

### Create a Box
```bash
raptor create box --name my-box
```

### Add a Secret to a Box
```bash
raptor create secret --box my-box --name API_KEY --value "sk_live_xxx"
```

### Generate a Random Password
```bash
raptor create password --length 32 --symbols
```

### List Secrets in a Box
```bash
raptor list secret --box my-box
```

### Get a Secret and Copy to Clipboard
```bash
raptor get secret --box my-box --name API_KEY --clip
```

### Edit a Secret
```bash
raptor edit secret --box my-box --name API_KEY
```

---

## Environment Variables

Raptor behavior can be customized with environment variables:

- **`CRYPTEX_FOLDER`**  
  Folder where Raptor stores boxes.  
  Default: `$HOME/.cryptex/boxes`

- **`CRYPTEX_BOX`**  
  Default box to use if `--box` is not provided.  

- **`RAPTOR_LOGLEVEL`**  
  Logging level: `debug`, `info`, `warn`, `error`  
  Default: `error`

- **`RAPTOR_TIMEOUT_SEC`**  
  Timeout in seconds of inactivity before Raptor exits.  
  Default: `600` (10 minutes)

---

## How It Works

- **Encryption**: Files and boxes are encrypted using strong, authenticated encryption. Each box is a JSON file stored in encrypted form.  
- **Secrets**: Inside a box, secrets are stored as key-value pairs. You can add, edit, list, and remove them without exposing other secrets.  
- **Passphrases**: Boxes are protected by passphrases. Raptor derives keys from passphrases securely (using a memory-hard KDF).  
- **Clipboard integration**: Secrets can be copied directly to clipboard, reducing accidental leaks in terminals.  

---

## Development

To build:

```bash
make build
```

To run tests:

```bash
make test
```

Lint and vet:

```bash
go vet ./...
golangci-lint run
```

---

## Contributing

Contributions are welcome!  
- Open issues for bugs or feature requests.  
- Submit PRs with tests and docs updated.  

---

## License

See [LICENSE](LICENSE).
