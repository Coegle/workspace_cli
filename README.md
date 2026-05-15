# Workspace CLI (`ws`)

[English](README.md) | [中文](README_zh.md)

`ws` is a command-line tool designed to solve the pain points of multi-repository microservice development. By leveraging Git Worktrees, `ws` helps you effortlessly manage multiple service repositories under isolated, feature-specific workspaces.

No more constantly stashing changes, checking out branches, or dealing with mixed-up local environments when switching between different tasks.

## 🚀 Features

- **Isolated Feature Workspaces:** Group multiple related microservices into a single directory based on a feature branch.
- **Git Worktree Powered:** Fast and lightweight. No need to clone repositories multiple times.
- **Interactive Service Resolution:** Fuzzy search for service names; if multiple match, interactively pick the one you need.
- **Auto-Sync:** Easily synchronize your local feature branches with their remote tracking branches.
- **Self-Updating:** Built-in update checker to ensure you are always using the latest version.
- **Cross-Platform:** Works on macOS, Linux, and Windows.

## 📦 Installation

### Option 1: Download Pre-compiled Binary (Recommended)
You can download the latest pre-compiled binaries for your operating system from the [Releases page](https://github.com/coegle/workspace_cli/releases).

### Option 2: Using `go install`
If you have Go installed (1.20+), you can install `ws` directly:
```bash
go install github.com/coegle/workspace_cli/cmd/ws@latest
```

## ⚙️ Configuration

Before using `ws`, you need to tell it where your base repositories are located and where you want to create your feature workspaces.

```bash
# Set the directory where your main git clones live
ws config base_repos /path/to/your/main/repos

# Set the directory where new feature workspaces should be created
ws config base_ws /path/to/your/workspaces
```

*By default, the configuration is stored in `~/.ws/config.yaml`.*

For more advanced configuration (like automatically setting up `symlinks` for your IDE across workspaces), please see the detailed examples in [`config.example.yaml`](config.example.yaml).

## 🛠️ Core Workflow

### 1. Create a Workspace & Add Services
Imagine you are working on a new feature called `feat_login` that requires modifying the `user_service` and `auth_service`.

```bash
# Syntax: ws add [-b base_branch] <feature_branch> <service1> <service2> ...
ws add feat_login user_service auth_service
```
This command will:
1. Create a new directory `feat_login` in your `base_ws` path.
2. Create a git worktree for `user_service` on the branch `feat_login`.
3. Create a git worktree for `auth_service` on the branch `feat_login`.

*Tip: If you want to branch off from a specific branch (e.g., `origin/develop`), use the `-b` flag: `ws add -b origin/develop feat_login user_service`.*

### 2. Open Workspace in IDE
Quickly open the entire feature workspace in your IDE (currently supports GoLand):
```bash
ws open feat_login
```

### 3. Sync Branches
When the base branch (e.g., `master`) has updates, you can quickly fast-forward your local branches (if they haven't diverged):
```bash
ws sync user_service
```

### 4. Remove a Service from Workspace
Done with `user_service`? Remove it from the workspace without affecting the main repository:
```bash
ws remove feat_login user_service
```

### 5. Destroy the Workspace
Once the feature is merged and done, clean up the entire workspace:
```bash
ws destroy feat_login
```
*(If there are uncommitted changes, `ws` will prompt you before destroying).*

## 🔄 Self Update

`ws` checks for updates periodically. If a new version is available, it will prompt you.
You can also trigger an update manually:
```bash
ws update
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
