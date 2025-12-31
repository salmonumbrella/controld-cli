# ControlD CLI Design

A command-line interface for the ControlD DNS management API, built on top of [controld-go](https://github.com/baptistecdr/controld-go).

## Goals

- Full ControlD API coverage for agent/automation use
- Design patterns matching airwallex-cli
- Keyring + environment variable credential storage
- Agent-friendly flags (--yes, --output json, --limit, --sort-by)

## Project Structure

```
controld-cli/
├── cmd/controld/main.go
├── internal/
│   ├── api/
│   │   ├── client.go               # Wraps controld-go with credential loading
│   │   ├── access.go               # Known IPs
│   │   ├── devices.go              # Devices + types
│   │   ├── profiles.go             # Profiles + options
│   │   ├── rules.go                # Custom rules
│   │   ├── default_rule.go         # Default rule
│   │   ├── filters.go              # Native + external filters
│   │   ├── folders.go              # Rule folders (groups)
│   │   ├── services.go             # Profile services
│   │   ├── misc.go                 # ListIP, ListNetwork
│   │   └── analytics.go            # Log levels, storage regions
│   ├── cmd/
│   │   ├── root.go                 # Root command + global flags
│   │   ├── auth.go                 # auth login/logout/list/status
│   │   ├── devices.go              # devices list/get/create/update/delete/types
│   │   ├── profiles.go             # profiles list/get/create/update/delete
│   │   ├── profiles_options.go     # profiles options list/set
│   │   ├── profiles_rules.go       # profiles rules list/create/update/delete
│   │   ├── profiles_default.go     # profiles default get/set
│   │   ├── profiles_filters.go     # profiles filters list/enable/disable
│   │   ├── profiles_folders.go     # profiles folders list/create/update/delete
│   │   ├── profiles_services.go    # profiles services list/set
│   │   ├── access.go               # access list/add/delete
│   │   ├── services.go             # services list/categories (reference)
│   │   ├── network.go              # network status/ip
│   │   ├── analytics.go            # analytics levels/regions (reference)
│   │   ├── users.go                # users (account info)
│   │   ├── version.go
│   │   └── completion.go
│   ├── config/paths.go             # App name, config paths
│   ├── secrets/store.go            # Keyring credential storage
│   ├── outfmt/                     # Output formatting (text/json)
│   ├── ui/                         # Colored terminal output
│   └── debug/                      # Debug logging
├── go.mod
└── go.sum
```

## Dependencies

```go
require (
    github.com/baptistecdr/controld-go  // ControlD API wrapper
    github.com/spf13/cobra              // CLI framework
    github.com/99designs/keyring        // Secure credential storage
)
```

## API Client Strategy

Thin wrapper approach - re-export controld-go's `*API` type, adding credential loading:

```go
func NewFromConfig(ctx context.Context) (*controld.API, error) {
    token := getToken() // keyring or env
    return controld.New(token)
}
```

Token resolution order:
1. `--token` flag
2. `CONTROLD_API_TOKEN` env var
3. `--account` flag → keyring lookup
4. Single account in keyring → auto-select
5. Multiple accounts → error with list

## Global Flags

```go
type rootFlags struct {
    Token  string  // --token (overrides all)
    Output string  // --output text|json
    Color  string  // --color auto|always|never
    Debug  bool    // --debug
    Query  string  // --query (jq filter)
    Yes    bool    // --yes/-y (skip prompts)
    Limit  int     // --limit (agent-friendly)
    SortBy string  // --sort-by
    Desc   bool    // --desc
}
```

Environment variables:
- `CONTROLD_API_TOKEN` - API token
- `CONTROLD_OUTPUT` - Default output format
- `CONTROLD_COLOR` - Color preference

## Command Tree

```
controld
├── auth login|logout|list|status
├── devices list|get|create|update|delete|types
├── profiles list|get|create|update|delete
│   ├── options list|set
│   ├── rules list|create|update|delete
│   ├── default get|set
│   ├── filters list|enable|disable
│   ├── folders list|create|update|delete
│   └── services list|set
├── access list|add|delete
├── services list|categories
├── network status|ip
├── analytics levels|regions
├── users
├── version
└── completion bash|zsh|fish
```

## Output Modes

```bash
# Text mode (default) - human readable
$ controld devices list
DEVICE_ID     NAME           STATUS    PROFILE
abc123        Home Router    active    Default

# JSON mode - agent/scripting friendly
$ controld devices list --output json
[{"device_id":"abc123","name":"Home Router",...}]

# JQ filtering
$ controld devices list --output json --query '.[].name'
```

## Auth Commands

```bash
# Interactive login
$ controld auth login
API Token: ****
✓ Authenticated as user@example.com

# Named accounts
$ controld auth login --name work --token cd_xxxx

# List accounts
$ controld auth list

# Check status
$ controld auth status

# Logout
$ controld auth logout [--name work]
```

## API Coverage

| Resource | Methods |
|----------|---------|
| Access | ListKnownIPs, LearnNewIPs, DeleteLearnedIPs |
| Account | ListUser |
| Analytics | ListLogLevels, ListStorageRegions |
| Devices | ListDevices, CreateDevice, UpdateDevice, DeleteDevice, ListDeviceType |
| Misc | ListIP, ListNetwork |
| Profiles | ListProfiles, CreateProfile, UpdateProfile, DeleteProfile, ListProfilesOptions, UpdateProfilesOption |
| Profile Rules | ListProfileCustomRules, CreateProfileCustomRule, UpdateProfileCustomRule, DeleteProfileCustomRule |
| Profile Default | ListProfileDefaultRule, UpdateProfileDefaultRule |
| Profile Filters | ListProfileNativeFilters, ListProfileExternalFilters, UpdateProfileFilter |
| Profile Folders | ListProfileRuleFolders, CreateProfileRuleFolder, UpdateProfileRuleFolder, DeleteProfileRuleFolder |
| Profile Services | ListProfileServices, UpdateProfileService |
| Services | ListServiceCategories, ListServices |

## Implementation Order

1. Scaffold - go.mod, main.go, root command, config
2. Auth - keyring store, login/logout/status commands
3. Core commands - devices, profiles
4. Profile sub-commands - rules, filters, services, folders
5. Reference commands - services, analytics, network
6. Polish - completion, version, upgrade check

## Repository

GitHub: `salmonumbrella/controld-cli`
