# Git Switch 使用指南

[English](USAGE.md)

## Profile 管理

### 创建 profile

```bash
git-switch profile add work
```

### 编辑 profile

```bash
git-switch profile edit work
```

使用默认编辑器（`$EDITOR`）打开 profile 文件。

### 列出所有 profiles

```bash
git-switch profile list
```

### 查看 profile 详情

```bash
git-switch profile show work
```

### 显示当前激活的 profile

```bash
git-switch profile current
```

### 为当前目录设置活跃 profile

```bash
git-switch profile use work
```

## 状态查看

显示当前激活的 profile 及匹配原因：

```bash
git-switch status
```

使用 `--quiet` 标志静默应用匹配的 profile，不输出内容（由 shell hook 内部调用）：

```bash
git-switch status --quiet
```

## 规则管理

### 列出所有匹配规则

```bash
git-switch rule list
```

### 添加匹配规则

```bash
git-switch rule add
```

### 删除匹配规则

```bash
git-switch rule remove
```

### 测试目录匹配

测试某个目录会匹配到哪个 profile：

```bash
git-switch rule test ~/projects/work/repo
```

## SSH 配置

重新生成 SSH config.d/git-switch：

```bash
git-switch ssh config
```

## Git 别名

### 安装

安装 `git sw` 别名、shell hook 和 shell 补全：

```bash
git-switch hook install
```

如果无法识别当前 shell，可以显式指定。支持的 shell 包括 `bash`、`zsh`、`powershell` 和 `pwsh`：

```bash
git-switch hook install --shell zsh
```

安装后可直接使用：

```bash
git sw status
git sw profile list
```

Shell hook 会在 `cd` 进入项目时自动切换 Git 身份。补全脚本提供所有 `git-switch` 命令的 tab 补全。

执行 `source ~/.zshrc`（或 `~/.bashrc`）或打开新终端即可在当前会话中启用补全。

### 查看状态

```bash
git-switch hook status
```

显示以下各项的安装状态：
- Git 别名（`git sw`）
- Shell hook（bash, zsh, powershell）
- Shell 补全（bash, zsh, powershell）

### 卸载

移除 `git sw` 别名、shell hook 和 shell 补全：

```bash
git-switch hook uninstall
```

如果无法识别当前 shell，可以显式指定：

```bash
git-switch hook uninstall --shell zsh
```

## Shell 补全

### 生成补全脚本

输出补全脚本到 stdout，供手动集成：

```bash
git-switch completion bash      # bash
git-switch completion zsh       # zsh
git-switch completion pwsh      # PowerShell
```

### 手动安装

如果不想使用 `hook install`，可以手动配置补全：

**bash** — 添加到 `~/.bashrc`：
```bash
source <(git-switch completion bash)
```

**zsh** — 添加到 `~/.zshrc`：
```zsh
source <(git-switch completion zsh)
```

**PowerShell** — 添加到 `$PROFILE`：
```powershell
git-switch completion pwsh | Out-String | Invoke-Expression
```

注：`hook install` 会自动完成以上步骤 —— 写入脚本、注入 RC 文件、为 zsh 设置 fpath。

## 项目模板

### 列出模板

```bash
git-switch template list
```

### 创建模板

```bash
git-switch template create
```

### 应用模板

```bash
git-switch template apply
```

## 卸载

```bash
# 完全卸载（删除配置）
git-switch uninstall

# 仅卸载程序，保留配置
git-switch uninstall --keep-config
```

## 版本查看

```bash
git-switch --version
```

通过 `go install` 安装时，版本号包含短 Git 提交哈希（如 `dev-a1b2c3d`）。官方发布二进制包含正式版本号。
