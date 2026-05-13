# Git Switch Usage Guide

[中文版](USAGE_zh.md)

## Profiles

### Create a profile

```bash
git-switch profile add work
```

### Edit a profile

```bash
vim ~/.config/git-switch/profiles/work.yaml
```

### List profiles

```bash
git-switch profile list
```

### Show profile details

```bash
git-switch profile show work
```

## Status

Show which profile is active and why:

```bash
git-switch status
```

## Rules

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

Install the `git sw` alias:

```bash
git-switch hook install
```

Then use:

```bash
git sw status
git sw profile list
```

## Templates

List available templates:

```bash
git-switch template list
```

## Uninstall

```bash
# Remove everything
git-switch uninstall

# Keep configuration
git-switch uninstall --keep-config
```
