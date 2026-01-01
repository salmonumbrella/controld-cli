# 🛡️ Control D CLI - DNS from the terminal. 

A powerful command-line interface for managing Control D - filter, devices, profiles, rules, and more.

## Features

- **Multiple account support** - manage multiple ControlD accounts
- **Devices** - create, configure, and manage DNS resolver endpoints
- **Profiles** - create and manage DNS filtering profiles with custom rules
- **Filters & Services** - enable/disable native filters and block/allow services
- **Network diagnostics** - test DNS resolution and view resolver status
- **Secure credential storage** using OS keyring (Keychain on macOS, Secret Service on Linux, Credential Manager on Windows)

## Installation

### Homebrew

```bash
brew install salmonumbrella/tap/controld-cli
```

or

```bash
go install github.com/salmonumbrella/controld-cli/cmd/controld@latest
```

## Quick Start

### 1. Authenticate

```bash
controld auth login
```

You'll be prompted to enter your ControlD API token. Get one from the [ControlD dashboard](https://controld.com/dashboard).

### 2. Verify Authentication

```bash
controld auth status
```

### 3. List Your Devices

```bash
controld devices list
```

## Configuration

### Account Selection

Specify the account using either a flag or environment variable:

```bash
# Via flag
controld devices list --account work

# Via environment
export CONTROLD_API_TOKEN=your-api-token
controld devices list
```

### Environment Variables

- `CONTROLD_API_TOKEN` - API token (alternative to keyring storage)
- `CONTROLD_OUTPUT` - Output format: `text` (default) or `json`
- `CONTROLD_COLOR` - Color mode: `auto` (default), `always`, or `never`
- `NO_COLOR` - Set to any value to disable colors (standard convention)

### Credential Storage

Credentials are stored securely in your OS keyring:
- macOS: Keychain
- Linux: Secret Service (GNOME Keyring, KWallet)
- Windows: Credential Manager

## Commands

### Authentication

```bash
controld auth login [--name <name>]       # Store API token
controld auth logout [--name <name>]      # Remove stored credentials
controld auth list                        # List configured accounts
controld auth status                      # Show authentication status
```

### Devices

```bash
controld devices list                                      # List all devices
controld devices get <deviceId>                            # Get device details
controld devices create --name <n> --profile-id <id>       # Create new device
controld devices modify <deviceId> [--name <n>] [--profile-id <id>] [--status <s>]
controld devices delete <deviceId>                         # Delete device
controld devices types                                     # List device types
```

### Profiles

```bash
controld profiles list                                     # List all profiles
controld profiles get <profileId>                          # Get profile details
controld profiles create --name <name> [--clone-from <id>] # Create new profile
controld profiles modify <profileId> [--name <name>]       # Modify profile
controld profiles delete <profileId>                       # Delete profile
```

### Profile Rules

```bash
controld profiles rules folders <profileId>               # List rule folders
controld profiles rules list <profileId> [--folder <id>]  # List custom rules
controld profiles rules create <profileId> --hostname <h> --action <a>  # Create rule
controld profiles rules delete <profileId> <hostname>     # Delete rule
```

Actions: `block`, `bypass`, `spoof`, `redirect`

### Profile Filters

```bash
controld profiles filters list <profileId>                 # List filters
controld profiles filters enable <profileId> <filterId>    # Enable filter
controld profiles filters disable <profileId> <filterId>   # Disable filter
```

### Profile Services

```bash
controld profiles services list <profileId>                           # List services
controld profiles services set <profileId> <serviceId> --action <a>   # Set action
controld profiles services disable <profileId> <serviceId>            # Remove rule
```

Actions: `block`, `bypass`, `spoof`

### Services Reference

```bash
controld services list [--category <cat>]                  # List available services
controld services categories                               # List service categories
```

### Network

```bash
controld network status                                    # List DNS resolvers
controld network resolve --host <domain>                   # Test DNS resolution
```

### Access (Known IPs)

```bash
controld access list <deviceId>                            # List known IPs
controld access add <deviceId> <ip> [<ip>...]              # Add known IPs
controld access delete <deviceId> <ip> [<ip>...]           # Remove known IPs
```

### Users

```bash
controld users                                             # Show account info
```

### Shell Completion

```bash
controld completion bash                                   # Generate bash completions
controld completion zsh                                    # Generate zsh completions
controld completion fish                                   # Generate fish completions
controld completion powershell                             # Generate PowerShell completions
```

## Output Formats

### Text

Human-readable tables with colors and formatting:

```bash
$ controld devices list
DEVICE_ID           NAME              STATUS    PROFILE
abc123def456...     Home Router       active    Family Safe
xyz789ghi012...     Work Laptop       active    Productivity

$ controld profiles list
PROFILE_ID          NAME              UPDATED
p1a2b3c4...         Family Safe       2024-01-15
p5d6e7f8...         Productivity      2024-01-20
```

### JSON

Machine-readable output for scripting and automation:

```bash
$ controld devices list --output json
[
  {
    "device_id": "abc123def456",
    "name": "Home Router",
    "status": 1,
    "profile": {"name": "Family Safe", "pk": "p1a2b3c4"}
  }
]
```

Data goes to stdout, errors and prompts to stderr for clean piping.

## Examples

### Create a device with a profile

```bash
# First, create a profile
controld profiles create --name "Kids Safe"

# Then create a device using that profile
controld devices create \
  --name "Kids Tablet" \
  --profile-id <profileId> \
  --icon tablet
```

### Block social media on a profile

```bash
# List available services to find social media
controld services list --category social

# Block specific services
controld profiles services set <profileId> facebook --action block
controld profiles services set <profileId> instagram --action block
controld profiles services set <profileId> tiktok --action block
```

### Add custom blocking rules

```bash
# Block a specific domain
controld profiles rules create <profileId> \
  --hostname "ads.example.com" \
  --action block

# Block multiple domains
controld profiles rules create <profileId> \
  --hostname "malware1.example.com" \
  --hostname "malware2.example.com" \
  --action block
```

### Enable ad blocking filters

```bash
# List available filters
controld profiles filters list <profileId>

# Enable ad blocking
controld profiles filters enable <profileId> ads
controld profiles filters enable <profileId> malware
```

### Test DNS resolution

```bash
# Check if a domain is blocked
controld network resolve --host "blocked-site.com"

# View resolver status
controld network status
```

### Export configuration as JSON

```bash
# Export all devices
controld devices list --output json > devices.json

# Export profile rules for backup
controld profiles rules list <profileId> --output json > rules-backup.json
```

### Switch between accounts

```bash
# Add multiple accounts
controld auth login --name personal
controld auth login --name work

# Use specific account
controld devices list --account personal
controld devices list --account work

# List configured accounts
controld auth list
```

## Global Flags

All commands support these flags:

- `--token <token>` - API token (overrides keyring and environment)
- `--account <name>` - Account name to use from keyring
- `--output <format>` - Output format: `text` or `json` (default: text)
- `--color <mode>` - Color mode: `auto`, `always`, or `never` (default: auto)
- `--yes` - Skip confirmation prompts
- `--debug` - Enable debug logging
- `--help` - Show help for any command

## Development

After cloning, install git hooks:

```bash
lefthook install
```

This installs [lefthook](https://github.com/evilmartians/lefthook) pre-commit and pre-push hooks for linting and testing.

Build locally:

```bash
go build -o controld ./cmd/controld
```

Run lints:

```bash
golangci-lint run
```

## License

MIT

## Links

- [ControlD Website](https://controld.com)
- [ControlD API Documentation](https://docs.controld.com/reference)