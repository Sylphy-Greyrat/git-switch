# git-switch

A cross-platform CLI tool for managing multiple Git user profiles and SSH keys.

## Features

- Multiple Git identities per project or remote URL
- SSH key configuration generation
- YAML configuration with nginx-style include support
- Git alias integration via `git sw`
- Cross-platform support for Linux, macOS, and Windows
- Project templates for quick setup

## Installation

```bash
go install github.com/sylphy/git-switch/cli@latest
```

## Quick Start

```bash
# Initialize configuration
git-switch init

# Edit your profile
vim ~/.config/git-switch/profiles/personal.yaml

# Install git alias
git-switch hook install

# Use it!
git sw status
git sw profile list
```

## Commands

| Command | Description |
|---------|-------------|
| `init` | Initialize configuration directory |
| `profile list` | List all profiles |
| `profile show <name>` | Show profile details |
| `profile add <name>` | Add a new profile |
| `profile remove <name>` | Remove a profile |
| `status` | Show current configuration status |
| `rule test <path>` | Test directory matching |
| `ssh config` | Regenerate SSH config |
| `hook install` | Install `git sw` alias |
| `hook uninstall` | Remove `git sw` alias |
| `template list` | List project templates |
| `uninstall` | Uninstall git-switch |

## Configuration

Configuration is stored in `~/.config/git-switch/`:

```text
~/.config/git-switch/
├── config.yaml           # Main configuration
└── profiles/
    ├── personal.yaml     # Personal profile
    └── work.yaml         # Work profile
```

See [docs/USAGE.md](docs/USAGE.md) for detailed usage instructions.

## License

MIT
