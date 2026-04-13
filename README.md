# CodeForge Resource Manager

CodeForge Resource Manager is a terminal UI (TUI) tool for defining and managing resource contracts, binding them to API operations, and validating alignment between resource definitions, field-level permissions, and API contract schemas. It helps produce consistent resource contracts for use with CodeForge Observer.

## Why This Exists

It helps ensure consistency between your API contracts and underlying data access rules by enforcing field-level permissions and resource bindings.

## Features

- Load resource contract and API contract files (JSON/OpenAPI) into app state
- Browse, add, and delete resources
- Edit field-level permissions (read, write, mutable) and persist changes
- Browse API endpoints and view bindings (via `x-resource`)
- Add/remove resource bindings on API operations
- Help view with global and page-specific keybindings
- TUI built with Bubble Tea / Lipgloss

## Installation

### Download (recommended)

Download the latest release from GitHub:
[Releases](https://github.com/Iztuk/codeforge-resource-manager/releases "CodeForge Resource Manager Releases")

Then:

```bash
chmod +x cfrm-linux-amd64
mv cfrm-linux-amd64 /usr/local/bin/cfrm
```

## Build From Source

```bash
git clone https://github.com/Iztuk/codeforge-resource-manager.git
cd codeforge-resource-manager
go build -o cfrm ./cmd
```

## Prerequisites

- Go 1.25.x (see `go.mod`)
- C compiler for sqlite3 (used by github.com/mattn/go-sqlite3)

On Ubuntu:

```bash
sudo apt-get install gcc libsqlite3-dev
```

## Run

- Example:
  `cfrm -api=/path/to/api_contract.json -resource=/path/to/resource_contract.json`
- Or with go run:
  `go run ./cmd/main.go -api=api_contract.json -resource=resource_contract.json`
- Both --api and --resource flags are required. The binary prints usage if flags
  are missing.

## TUI Controls (default)

- Navigation: `h / j / k / l`
- Select / open: `enter`
- Back: `backspace`
- Exit: `ctrl+c`
- Add: `ctrl+a`
- Remove/Delete: `ctrl+d`
- Toggle Field Permission: `space`

## File Formats and Expectations

- Resource contract: JSON describing resources/tables and per-field permissions.
  The app reads and writes this file to persist resource changes and permission updates.

- API contract: OpenAPI file. The app reads endpoints and stores bindings using an `x-resource` extension on operations to link them resource definitions.

## Development Notes

- TUI implemented using charmbracelet libraries
- State initialization happens via `internal/state.InitializeAppState(apiPath, resourcePath)`
- Entry point: cmd/main.go (flags: `-api`, `-resource`)

## Roadmap

- Cross-platform releases (macOS / Windows)

## Contributing

Contributions are welcome.

If you plan to make significant changes, please open an issue to discuss the approach.

## License

This project is licensed under the GNU General Public License v3 (GPLv3).

See the LICENSE file for details.
