# git-switch

跨平台 CLI 工具，用于管理多个 Git 用户配置和 SSH 密钥。

[English](README.md)

## 功能特性

- 按项目或远程 URL 自动切换 Git 身份
- SSH 密钥配置自动生成
- YAML 配置，支持 nginx 风格的 include 机制
- Git 别名集成（`git sw`）
- 跨平台支持：Linux、macOS、Windows
- 项目模板，快速初始化

## 安装

### 二进制下载

从 [GitHub Releases](https://github.com/Sylphy-Greyrat/git-switch/releases/latest) 下载：

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

### Go 安装

```bash
go install github.com/sylphy/git-switch/cli@latest
```

## 快速开始

```bash
# 初始化配置
git-switch init

# 编辑你的 profile
vim ~/.config/git-switch/profiles/personal.yaml

# 安装 git 别名
git-switch hook install

# 开始使用！
git sw status
git sw profile list
```

## 命令一览

| 命令 | 说明 |
|------|------|
| `init` | 初始化配置目录 |
| `profile list` | 列出所有 profiles |
| `profile show <name>` | 查看 profile 详情 |
| `profile add <name>` | 添加新 profile |
| `profile remove <name>` | 删除 profile |
| `profile current` | 显示当前激活的 profile |
| `profile edit <name>` | 在默认编辑器中编辑 profile |
| `profile use <name>` | 为当前目录设置活跃 profile |
| `status` | 显示当前配置状态 |
| `rule list` | 列出所有匹配规则 |
| `rule add` | 添加匹配规则 |
| `rule remove` | 删除匹配规则 |
| `rule test <path>` | 测试目录匹配规则 |
| `ssh config` | 重新生成 SSH 配置 |
| `hook install` | 安装 `git sw` 别名 |
| `hook uninstall` | 移除 `git sw` 别名 |
| `hook status` | 显示 hook 安装状态 |
| `template list` | 列出项目模板 |
| `template create` | 创建新项目模板 |
| `template apply` | 将模板应用到项目目录 |
| `uninstall` | 卸载 git-switch |
| `--version` | 显示版本 |

## 配置说明

配置存储在 `~/.config/git-switch/` 目录下：

```text
~/.config/git-switch/
├── config.yaml           # 主配置文件
└── profiles/
    ├── personal.yaml     # 个人 profile
    └── work.yaml         # 工作 profile
```

详细用法请参阅 [docs/USAGE_zh.md](docs/USAGE_zh.md)。

## 许可证

MIT
