# Workspace CLI (`ws`)

[English](README.md) | [中文](README_zh.md)

`ws` 是一个旨在解决多仓库微服务开发痛点的命令行工具。通过充分利用 Git Worktrees 的能力，`ws` 帮助你轻松地在独立的“功能工作区”中管理多个相关的服务仓库。

告别在不同任务间切换时频繁地 stash 代码、切换分支以及本地环境错乱的烦恼！

## 🚀 核心特性

- **隔离的功能工作区：** 基于功能分支，将多个相关的微服务组合到一个独立的物理目录中。
- **基于 Git Worktree：** 快速且轻量。无需多次 clone 同一个仓库。
- **交互式服务解析：** 支持服务名称模糊搜索；如果有多个匹配项，支持交互式选择。
- **分支自动同步：** 轻松将本地功能分支与远程跟踪分支（或主干分支）进行同步。
- **CoW 缓存预热（macOS）：** 将主仓库中被 git 忽略的 codegen 目录（如 Kitex 的 `kitex_gen`）通过写时复制克隆到新工作区——近乎零耗时、零额外占盘。
- **生命周期钩子（Hooks）：** 通过 `~/.ws/hooks/` 下的简单约定，在 `ws add` 前后运行你自己的脚本（如拉取最新主干分支、安装依赖、自动打开 IDE）。
- **内置自更新：** 内置版本检查机制，确保你使用的始终是最新版本。
- **跨平台：** 完美支持 macOS, Linux 和 Windows。

## 📦 安装指南

### 方式 1: 下载预编译二进制文件 (推荐)
你可以从 [Releases 页面](https://github.com/coegle/workspace_cli/releases) 下载适用于你操作系统的最新版预编译程序。

### 方式 2: 使用 `go install`
如果你已经安装了 Go (1.20+)，可以直接通过以下命令安装：
```bash
go install github.com/coegle/workspace_cli/cmd/ws@latest
```

## ⚙️ 初始配置

### 必选配置

在使用 `ws` 之前，你需要告诉它你的基础代码仓库在哪里，以及你想把新的功能工作区建在哪里。

```bash
# 设置你主干代码库（git clones）所在的目录
ws config base_repos /path/to/your/main/repos

# 设置你希望新建功能工作区的目标目录
ws config base_ws /path/to/your/workspaces
```

*配置默认保存在 `~/.ws/config.yaml` 文件中。*

### 可选配置

以下配置均为可选——只需上面两项，`ws` 即可开箱即用。完整配置项（如 `replace_slash`、用于跨工作区共享 IDE 配置的 `symlinks`，以及下面的 `cow_dirs`）请参考 [`config.example.yaml`](config.example.yaml) 示例文件。

#### CoW 缓存预热（macOS / APFS）

有些仓库会保存体积较大的**被 git 忽略的 codegen 产物**（例如 Kitex 的 `kitex_gen`、Hertz 的 `hertz_gen`，或 protobuf/gRPC 生成的 `gen/` 目录）。由于这些目录被忽略，`git worktree add` 不会 checkout 它们，于是每个新工作区都得重新生成一遍——既慢又占盘。

通过 `cow_dirs`，`ws` 会用 APFS 的 `clonefile` 把这些目录从**主仓库**写时复制（CoW）克隆到每个新工作区。数据块与主仓库那份共享，因此近乎零耗时、几乎不额外占盘。

```yaml
# 位于 ~/.ws/config.yaml
cow_dirs:
  - kitex_gen
```

注意事项：
- Donor（克隆源）永远是主仓库（`base_repos/<service>/kitex_gen`），所以请确保你至少在主仓库里跑过一次 codegen。
- 它是纯缓存/预热：原样克隆，**不**校验 IDL 一致性。如果新工作区的 IDL 有变化，在工作区里重新跑一次 codegen 即可。
- 如果 donor 目录不存在，该条目会被静默跳过。
- 仅在同一 APFS 卷内生效；跨卷 / 非 APFS 会直接跳过（不会退化成完整拷贝）。

#### 生命周期钩子（Hooks）

`ws` 可以在 `ws add` 前后运行你自己的脚本，让你无需修改 `ws` 本身即可定制工作流（例如在创建 worktree 前拉取最新主干分支、安装依赖，或在完成后自动打开 IDE）。

Hooks 采用**按约定开启**：只需在约定路径放一个可执行脚本，`ws` 就会运行它。如果脚本不存在，则静默跳过。

| Hook | 触发时机 | 失败行为 |
| --- | --- | --- |
| `~/.ws/hooks/pre-add` | 工作区创建之前 | **中断** `ws add`（脚本非 0 退出时） |
| `~/.ws/hooks/post-add` | 所有 worktree 就绪之后 | 尽力而为；仅打印警告，`ws` 仍视为成功 |

每个脚本通过环境变量获得工作区上下文：

| 变量 | 含义 |
| --- | --- |
| `WS_DIR` | 工作区目录的绝对路径 |
| `WS_NAME` | 工作区（安全化后的目录）名称 |
| `WS_BRANCH` | 功能分支名 |
| `WS_REPOS` | 存放源仓库的 `base_repos` 路径 |

示例——创建前校验分支名，完成后打开 IDE：

```bash
# ~/.ws/hooks/pre-add （记得 chmod +x）
#!/usr/bin/env bash
[[ "$WS_BRANCH" == feat/* ]] || { echo "分支名必须以 feat/ 开头"; exit 1; }
```

```bash
# ~/.ws/hooks/post-add （chmod +x）
#!/usr/bin/env bash
# 此时工作目录已是 $WS_DIR
code .
```

注意事项：
- 脚本必须可执行（`chmod +x`）；不可执行的文件会被跳过并打印警告。
- `post-add` 运行时工作目录为 `$WS_DIR`；`pre-add` 运行时该目录尚未创建，因此使用当前目录。



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
