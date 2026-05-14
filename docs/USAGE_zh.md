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

安装 `git sw` 别名，并自动为当前 shell 安装 hook：

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

### 查看状态

```bash
git-switch hook status
```

### 卸载别名

```bash
git-switch hook uninstall
```

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
