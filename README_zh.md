# Workspace CLI (`ws`)

[English](README.md) | [中文](README_zh.md)

`ws` 是一个旨在解决多仓库微服务开发痛点的命令行工具。通过充分利用 Git Worktrees 的能力，`ws` 帮助你轻松地在独立的“功能工作区”中管理多个相关的服务仓库。

告别在不同任务间切换时频繁地 stash 代码、切换分支以及本地环境错乱的烦恼！

## 🚀 核心特性

- **隔离的功能工作区：** 基于功能分支，将多个相关的微服务组合到一个独立的物理目录中。
- **基于 Git Worktree：** 快速且轻量。无需多次 clone 同一个仓库。
- **交互式服务解析：** 支持服务名称模糊搜索；如果有多个匹配项，支持交互式选择。
- **分支自动同步：** 轻松将本地功能分支与远程跟踪分支（或主干分支）进行同步。
- **内置自更新：** 内置版本检查机制，确保你使用的始终是最新版本。
- **跨平台：** 完美支持 macOS, Linux 和 Windows。

## 📦 安装指南

### 方式 1: 下载预编译二进制文件 (推荐)
你可以从 [Releases 页面](https://github.com/coegle/workspace_cli/releases) 下载适用于你操作系统的最新版预编译程序。

### 方式 2: 使用 `go install`
如果你已经安装了 Go (1.20+)，可以直接通过以下命令安装：
```bash
go install github.com/coegle/workspace_cli@latest
# 注意：安装后如果名字不是 ws，你可以重命名或在 shell 中配置 alias。
```

## ⚙️ 初始配置

在使用 `ws` 之前，你需要告诉它你的基础代码仓库在哪里，以及你想把新的功能工作区建在哪里。

```bash
# 设置你主干代码库（git clones）所在的目录
ws config base_repos /path/to/your/main/repos

# 设置你希望新建功能工作区的目标目录
ws config base_ws /path/to/your/workspaces
```

*配置默认保存在 `~/.ws/config.yaml` 文件中。*

## 🛠️ 核心工作流

### 1. 创建工作区 & 添加服务
假设你正在开发一个名为 `feat_login` 的新功能，需要同时修改 `user_service` 和 `auth_service`。

```bash
# 语法: ws add [-b base_branch] <feature_branch> <service1> <service2> ...
ws add feat_login user_service auth_service
```
该命令会自动执行以下操作：
1. 在你的 `base_ws` 路径下创建一个名为 `feat_login` 的新目录。
2. 为 `user_service` 创建一个指向 `feat_login` 分支的 git worktree。
3. 为 `auth_service` 创建一个指向 `feat_login` 分支的 git worktree。

*提示：如果你想基于特定的分支（例如 `origin/develop`）切出新分支，可以使用 `-b` 参数：`ws add -b origin/develop feat_login user_service`。*

### 2. 在 IDE 中打开工作区
一键在 IDE 中打开整个功能工作区（目前原生支持 GoLand）：
```bash
ws open feat_login
```

### 3. 同步分支
当基础分支（如 `master`）有更新时，你可以快速通过 fast-forward 同步你的本地分支：
```bash
ws sync user_service
```

### 4. 从工作区移除服务
`user_service` 的开发完成了？将它从当前工作区移除，这不会影响你的主仓库：
```bash
ws remove feat_login user_service
```

### 5. 销毁整个工作区
功能开发完毕并合并后，一键清理整个工作区目录：
```bash
ws destroy feat_login
```
*(如果检测到有未提交的代码，`ws` 会在销毁前安全提示)。*

## 🔄 自动更新

`ws` 会定期在后台检查更新。如果有新版本，它会在终端中提醒你。
你也可以随时手动触发更新：
```bash
ws update
```

## 📄 开源协议

本项目基于 MIT 协议开源 - 查看 [LICENSE](LICENSE) 文件了解更多详情。
