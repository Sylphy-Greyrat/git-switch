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

Install the `git sw` alias:

```bash
git-switch hook install
```

Then use:

```bash
git sw status
git sw profile list
```

### Check status

```bash
git-switch hook status
```

### Uninstall

```bash
git-switch hook uninstall
```

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
