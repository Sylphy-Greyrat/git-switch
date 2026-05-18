# Git Switch Usage Guide

[中文版](USAGE_zh.md)

## Profiles

### Create a profile

```bash
git-switch profile add work
```

### Edit a profile

```bash
git-switch profile edit work
```

Opens the profile in your default editor (`$EDITOR`).

### List profiles

```bash
git-switch profile list
```

### Show profile details

```bash
git-switch profile show work
```

### Show current active profile

```bash
git-switch profile current
```

### Set active profile for current directory

```bash
git-switch profile use work
```

## Status

Show which profile is active and why:

```bash
git-switch status
```

Use `--quiet` flag to silently apply the matched profile without output (used internally by shell hook):

```bash
git-switch status --quiet
```

## Rules

### List all matching rules

```bash
git-switch rule list
```

### Add a matching rule

```bash
git-switch rule add
```

### Remove a matching rule

```bash
git-switch rule remove
```

### Test directory matching

Test which profile would match a directory:

```bash
git-switch rule test ~/projects/work/repo
```

## SSH Configuration

Regenerate SSH config.d/git-switch:

```bash
git-switch ssh config
```

## Git Alias

### Install

Install the `git sw` alias, shell hook, and shell completion for the current shell:

```bash
git-switch hook install
```

If the current shell cannot be detected, specify it explicitly. Supported shells are `bash`, `zsh`, `powershell`, and `pwsh`:

```bash
git-switch hook install --shell zsh
```

After installation, you can use:

```bash
git sw status
git sw profile list
```

The shell hook automatically switches your Git identity when you `cd` into a project. The completion script provides tab-completion for all `git-switch` commands.

Reload your shell or run `source ~/.zshrc` (or `~/.bashrc`) for completion to take effect in the current session.

### Check status

```bash
git-switch hook status
```

Shows installation status for:
- Git alias (`git sw`)
- Shell hook (bash, zsh, powershell)
- Shell completion (bash, zsh, powershell)

### Uninstall

Remove the `git sw` alias, shell hook, and shell completion:

```bash
git-switch hook uninstall
```

If the current shell cannot be detected, specify it explicitly:

```bash
git-switch hook uninstall --shell zsh
```

## Shell Completion

### Generate completion script

Output the completion script to stdout for manual integration:

```bash
git-switch completion bash      # bash
git-switch completion zsh       # zsh
git-switch completion pwsh      # PowerShell
```

### Manual installation

If you prefer not to use `hook install`, you can manually configure completion:

**bash** — add to `~/.bashrc`:
```bash
source <(git-switch completion bash)
```

**zsh** — add to `~/.zshrc`:
```zsh
source <(git-switch completion zsh)
```

**PowerShell** — add to `$PROFILE`:
```powershell
git-switch completion pwsh | Out-String | Invoke-Expression
```

Note: `hook install` handles all of this automatically — writing scripts, injecting RC files, and setting up fpath for zsh.

## Templates

### List templates

```bash
git-switch template list
```

### Create a template

```bash
git-switch template create
```

### Apply a template

```bash
git-switch template apply
```

## Uninstall

```bash
# Remove everything
git-switch uninstall

# Keep configuration
git-switch uninstall --keep-config
```

## Version

```bash
git-switch --version
```

When installed via `go install`, the version includes a short Git commit hash (e.g., `dev-a1b2c3d`). Official release binaries include the release version instead.
