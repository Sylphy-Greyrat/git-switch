# git-switch

A cross-platform CLI tool for managing multiple Git user profiles and SSH keys.

[中文版](README_zh.md)

## Features

- Multiple Git identities per project or remote URL
- SSH key configuration generation
- YAML configuration with nginx-style include support
- Shell completion for bash, zsh, and PowerShell
- Git alias integration via `git sw` with cd-auto-switch hook
- Cross-platform support for Linux, macOS, and Windows
- Project templates for quick setup

## Installation

### Binary Download

Download from [GitHub Releases](https://github.com/Sylphy-Greyrat/git-switch/releases/latest):

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/Sylphy-Greyrat/git-switch/releases/download/v0.1.0/git-switch-v0.1.0-aarch64-macos
chmod +x git-switch-v0.1.0-aarch64-macos
sudo mv git-switch-v0.1.0-aarch64-macos /usr/local/bin/git-switch

# macOS (Intel)
curl -LO https://github.com/Sylphy-Greyrat/git-switch/releases/download/v0.1.0/git-switch-v0.1.0-x86_64-macos
chmod +x git-switch-v0.1.0-x86_64-macos
sudo mv git-switch-v0.1.0-x86_64-macos /usr/local/bin/git-switch

# Linux (x86_64)
curl -LO https://github.com/Sylphy-Greyrat/git-switch/releases/download/v0.1.0/git-switch-v0.1.0-x86_64-linux
chmod +x git-switch-v0.1.0-x86_64-linux
sudo mv git-switch-v0.1.0-x86_64-linux /usr/local/bin/git-switch
```

### Go Install

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
| `profile current` | Show current active profile |
| `profile edit <name>` | Edit profile in default editor |
| `profile use <name>` | Set active profile for current directory |
| `status` | Show current configuration status |
| `rule list` | List all matching rules |
| `rule add` | Add a matching rule to a profile |
| `rule remove` | Remove a matching rule |
| `rule test <path>` | Test directory matching |
| `ssh config` | Regenerate SSH config |
| `hook install` | Install `git sw` alias, shell hook, and completion |
| `hook uninstall` | Remove `git sw` alias, shell hook, and completion |
| `hook status` | Show hook and completion installation status |
| `completion <shell>` | Generate shell completion script |
| `template list` | List project templates |
| `template create` | Create a new project template |
| `template apply` | Apply template to project directory |
| `uninstall` | Uninstall git-switch |
| `--version` | Show version |

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
