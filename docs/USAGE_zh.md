# Git Switch 使用指南

[English](USAGE.md)

## Profile 管理

### 创建 profile

```bash
git-switch profile add work
```

### 编辑 profile

```bash
vim ~/.config/git-switch/profiles/work.yaml
```

### 列出所有 profiles

```bash
git-switch profile list
```

### 查看 profile 详情

```bash
git-switch profile show work
```

## 状态查看

显示当前激活的 profile 及匹配原因：

```bash
git-switch status
```

## 规则测试

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

安装 `git sw` 别名：

```bash
git-switch hook install
```

安装后可直接使用：

```bash
git sw status
git sw profile list
```

## 项目模板

列出可用模板：

```bash
git-switch template list
```

## 卸载

```bash
# 完全卸载（删除配置）
git-switch uninstall

# 仅卸载程序，保留配置
git-switch uninstall --keep-config
```
