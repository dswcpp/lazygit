<div align="center">
  <img width="536" src="https://user-images.githubusercontent.com/8456633/174470852-339b5011-5800-4bb9-a628-ff230aa8cd4e.png">
</div>

<div align="center">

# Lazygit - Git 命令行终端用户界面工具

[![GitHub Releases](https://img.shields.io/github/downloads/jesseduffield/lazygit/total)](https://github.com/jesseduffield/lazygit/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/jesseduffield/lazygit)](https://goreportcard.com/report/github.com/jesseduffield/lazygit) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/f46416b715d74622895657935fcada21)](https://app.codacy.com/gh/jesseduffield/lazygit/dashboard) [![golangci-lint](https://img.shields.io/badge/linted%20by-golangci--lint-brightgreen)](https://golangci-lint.run/) [![GitHub tag](https://img.shields.io/github/v/tag/jesseduffield/lazygit?color=blue)](https://github.com/jesseduffield/lazygit/releases/latest)

![commit_and_push](../assets/demo/commit_and_push-compressed.gif)

</div>

## 一、项目概述

Lazygit 是一款基于终端的 Git 命令简化工具，旨在为开发人员提供直观、高效的版本控制操作界面。本工具通过图形化终端界面，简化了 Git 的复杂操作流程，显著提升开发效率。

### 1.1 核心优势

传统 Git 命令行操作存在以下痛点：
- 交互式变基需要手动编辑 TODO 文件
- 部分暂存文件需要逐个处理代码块
- 分支切换时频繁遇到不必要的暂存提示

Lazygit 通过可视化界面和快捷键操作，彻底解决上述问题，让版本控制回归简单高效。

### 1.2 赞助支持

本项目的持续维护得益于全体[贡献者](https://github.com/jesseduffield/lazygit/graphs/contributors)和[赞助者](https://github.com/sponsors/jesseduffield)的支持。如需赞助本项目，请访问 [GitHub Sponsors](https://github.com/sponsors/jesseduffield)。

<p align="center">
<!-- sponsors --><a href="https://github.com/intabulas"><img src="https://github.com/intabulas.png" width="60px" alt="Mark Lussier" /></a><a href="https://github.com/peppy"><img src="https://github.com/peppy.png" width="60px" alt="Dean Herbert" /></a><a href="https://github.com/piot"><img src="https://github.com/piot.png" width="60px" alt="Peter Bjorklund" /></a><a href="https://github.com/rgwood"><img src="https://github.com/rgwood.png" width="60px" alt="Reilly Wood" /></a><a href="https://github.com/oliverguenther"><img src="https://github.com/oliverguenther.png" width="60px" alt="Oliver Günther" /></a><!-- 更多赞助者请访问项目主页 -->
</p>

## 二、功能特性

### 2.1 逐行暂存

在选定行上按空格键即可暂存该行，按 `v` 键可选择行范围，按 `a` 键可选择整个代码块。

![stage_lines](../assets/demo/stage_lines-compressed.gif)

### 2.2 交互式变基

按 `i` 键启动交互式变基。支持压缩提交 (`s`)、修正提交 (`f`)、删除提交 (`d`)、编辑提交 (`e`)、上移 (`ctrl+k`) 或下移 (`ctrl+j`) 提交。按 `m` 键打开变基选项菜单，选择"继续"完成变基操作。

![interactive_rebase](../assets/demo/interactive_rebase-compressed.gif)

### 2.3 拣选提交

在提交上按 `shift+c` 复制提交，按 `shift+v` 粘贴（拣选）提交。

![cherry_pick](../assets/demo/cherry_pick-compressed.gif)

### 2.4 二分查找

在提交视图中按 `b` 键，将提交标记为好/坏，以启动 git bisect 二分查找。

![bisect](../assets/demo/bisect-compressed.gif)

### 2.5 清空工作树

按 `shift+d` 打开重置选项菜单，选择"清空"选项，可清除 `git status` 显示的所有内容（包括脏子模块）。

![Nuke working tree](../assets/demo/nuke_working_tree-compressed.gif)

### 2.6 修正旧提交

在任意提交上按 `shift+a`，可使用当前暂存的更改修正该提交（后台运行交互式变基）。

![amend_old_commit](../assets/demo/amend_old_commit-compressed.gif)

### 2.7 过滤视图

按 `/` 键可过滤视图。例如，过滤分支视图后按回车键查看其提交。

![filter](../assets/demo/filter-compressed.gif)

### 2.8 自定义命令

Lazygit 提供灵活的[自定义命令系统](docs/Custom_Command_Keybindings.md)。可定义自定义命令模拟内置操作。

![custom_command](../assets/demo/custom_command-compressed.gif)

### 2.9 工作树管理

可创建工作树以同时处理多个分支，无需暂存或创建 WIP 提交。在分支视图中按 `w` 键，从选定分支创建工作树并切换。

![worktree_create_from_branches](../assets/demo/worktree_create_from_branches-compressed.gif)

### 2.10 变基魔法（自定义补丁）

可从旧提交构建自定义补丁，然后从提交中删除补丁、拆分新提交、反向应用补丁到索引等。

详见 [Rebase magic Youtube 教程](https://youtu.be/4XaToVut_hs)。

![custom_patch](../assets/demo/custom_patch-compressed.gif)

### 2.11 从标记基础提交变基

按 `shift+b` 标记基础提交，然后在目标分支上按 `r` 进行变基，仅带入功能分支的提交。

![rebase_onto](../assets/demo/rebase_onto-compressed.gif)

### 2.12 撤销/重做

按 `z` 键撤销上一操作，按 `ctrl+z` 重做。撤销功能使用 reflog，仅适用于提交和分支操作。

详见[撤销文档](/docs/Undoing.md)。

![undo](../assets/demo/undo-compressed.gif)

### 2.13 提交图

在放大窗口中查看提交图（使用 `+` 和 `_` 循环切换屏幕模式），提交图将显示。颜色对应提交作者，导航时高亮显示父提交。

![commit_graph](../assets/demo/commit_graph-compressed.gif)

### 2.14 比较两个提交

在提交（或分支/引用）上按 `shift+w`，打开菜单标记该提交，选择另一提交后将显示差异。按回车键查看差异文件。

![diff_commits](../assets/demo/diff_commits-compressed.gif)

## 三、教程资源

- [15 分钟掌握 15 个 Lazygit 功能](https://youtu.be/CPLdltN7wgE)
- [基础教程](https://youtu.be/VDXvbHZYeKY)
- [变基魔法教程](https://youtu.be/4XaToVut_hs)

## 四、安装指南

[![Packaging status](https://repology.org/badge/vertical-allrepos/lazygit.svg?columns=3)](https://repology.org/project/lazygit/versions)

*注：上述大部分软件包由第三方维护，请自行验证维护者可信度。*

### 4.1 二进制发行版

Windows、Mac OS (10.12+) 或 Linux 系统可下载二进制发行版：[下载地址](../../releases)

### 4.2 开发容器特性

如需在 GitHub Codespaces 中使用 lazygit，可使用基于二进制发行版的第三方[开发容器特性](https://github.com/GeorgOfenbeck/features/tree/main/src/lazygit-linuxbinary)。

### 4.3 Homebrew（支持 Linux）

```sh
brew install lazygit
```

### 4.4 MacPorts

```sh
sudo port install lazygit
```

### 4.5 Void Linux

```sh
sudo xbps-install -S lazygit
```

### 4.6 Scoop (Windows)

```sh
# 添加 extras bucket
scoop bucket add extras

# 安装 lazygit
scoop install lazygit
```

### 4.7 Arch Linux

- 稳定版：`sudo pacman -S lazygit`
- 开发版：<https://aur.archlinux.org/packages/lazygit-git/>

AUR 安装说明：<https://wiki.archlinux.org/index.php/Arch_User_Repository>

### 4.8 Fedora / Amazon Linux 2023 / CentOS Stream

```sh
sudo dnf copr enable dejan/lazygit
sudo dnf install lazygit
```

### 4.9 Debian 和 Ubuntu

**Debian 13 "Trixie"、Sid 及更高版本，或 Ubuntu 25.10 "Questing Quokka" 及更高版本：**

```sh
sudo apt install lazygit
```

**Debian 12 "Bookworm"、Ubuntu 25.04 "Plucky Puffin" 及更早版本：**

```sh
LAZYGIT_VERSION=$(curl -s "https://api.github.com/repos/jesseduffield/lazygit/releases/latest" | \grep -Po '"tag_name": *"v\K[^"]*')
curl -Lo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/download/v${LAZYGIT_VERSION}/lazygit_${LAZYGIT_VERSION}_Linux_x86_64.tar.gz"
tar xf lazygit.tar.gz lazygit
sudo install lazygit -D -t /usr/local/bin/
```

验证安装：

```sh
lazygit --version
```

### 4.10 NixOS

**使用 nixpkgs 中的 lazygit：**

```sh
nix-shell -p lazygit
# 或启用 flakes
nix run nixpkgs#lazygit
```

或将 lazygit 添加到 `configuration.nix` 的 `environment.systemPackages` 选项。

**使用官方 lazygit flake：**

```sh
# 直接从仓库运行
nix run github:jesseduffield/lazygit

# 从源码构建
nix build github:jesseduffield/lazygit

# 开发环境
nix develop github:jesseduffield/lazygit
```

### 4.11 FreeBSD

```sh
pkg install lazygit
```

### 4.12 Termux

```sh
apt install lazygit
```

### 4.13 Conda

```sh
conda install -c conda-forge lazygit
```

### 4.14 Go

```sh
go install github.com/jesseduffield/lazygit@latest
```

*注：如提示找不到 lazygit，需将 `~/go/bin` (MacOS/Linux) 或 `%HOME%\go\bin` (Windows) 添加到 $PATH。*

### 4.15 Chocolatey (Windows)

```sh
choco install lazygit
```

### 4.16 Winget (Windows 10 1709 或更高版本)

```powershell
winget install -e --id=JesseDuffield.lazygit
```

### 4.17 手动安装

需先[安装 Go](https://golang.org/doc/install)：

```sh
git clone https://github.com/jesseduffield/lazygit.git
cd lazygit
go install
```

或使用 `go run main.go` 一步编译运行。

## 五、使用说明

### 5.1 基本用法

在 git 仓库目录内执行：

```sh
$ lazygit
```

可添加别名简化调用：

```sh
echo "alias lg='lazygit'" >> ~/.zshrc
```

### 5.2 快捷键绑定

快捷键列表详见[快捷键文档](/docs/keybindings)。

### 5.3 退出时切换目录

如需在退出 lazygit 时切换到仓库目录，在 `~/.zshrc` 中添加：

```sh
lg()
{
    export LAZYGIT_NEW_DIR_FILE=~/.lazygit/newdir

    lazygit "$@"

    if [ -f $LAZYGIT_NEW_DIR_FILE ]; then
            cd "$(cat $LAZYGIT_NEW_DIR_FILE)"
            rm -f $LAZYGIT_NEW_DIR_FILE > /dev/null
    fi
}
```

执行 `source ~/.zshrc` 后，调用 `lg` 并退出时将自动切换目录。使用 `shift+Q` 退出可覆盖此行为。

### 5.4 撤销/重做

详见[撤销文档](/docs/Undoing.md)。

## 六、配置说明

### 6.1 配置文档

详见[配置文档](docs/Config.md)。

### 6.2 自定义分页器

详见[自定义分页器文档](docs/Custom_Pagers.md)。

### 6.3 自定义命令

如 lazygit 缺少某功能，可通过自定义命令实现。详见[自定义命令文档](docs/Custom_Command_Keybindings.md)。

### 6.4 Git flow 支持

Lazygit 支持 [Gitflow](https://github.com/nvie/gitflow)。了解 Gitflow 模型请参阅 Vincent Driessen 的[原始文章](https://nvie.com/posts/a-successful-git-branching-model/)。在分支视图中按 `i` 查看 Gitflow 选项。

## 七、贡献指南

### 7.1 参与贡献

欢迎您的参与！请查阅[贡献指南](CONTRIBUTING.md)。

如需讨论仓库外的贡献者话题，请加入 [Discord 频道](https://discord.gg/ehwFt2t4wt)。

<a href="https://discord.gg/ehwFt2t4wt"><img src='../assets/discord.png' width='75'></a>

观看此[视频](https://www.youtube.com/watch?v=kNavnhzZHtk)了解如何在 lazygit 中创建小功能。

### 7.2 本地调试

在一个终端标签页运行 `lazygit --debug`，在另一个标签页运行 `lazygit --logs`，可并排查看程序及其日志输出。

## 八、捐赠支持

如需支持 lazygit 的开发，请考虑[赞助我](https://github.com/sponsors/jesseduffield)（GitHub 将在 12 个月内按 1:1 匹配所有捐赠）。

## 九、常见问题

### 9.1 提交颜色代表什么？

- 绿色：提交已包含在 master 分支中
- 黄色：提交未包含在 master 分支中
- 红色：提交尚未推送到上游分支

## 十、相关链接

### 10.1 开发者信息

如需了解 Jesse 的开发动态，请关注 [Twitter](https://twitter.com/DuffieldJesse) 或访问[博客](https://jesseduffield.com/)。

### 10.2 替代工具

如 lazygit 不能完全满足需求，可考虑以下替代工具：

- [GitUI](https://github.com/Extrawurst/gitui)
- [tig](https://github.com/jonas/tig)
- [GitArbor TUI](https://github.com/cadamsdev/gitarbor-tui)

---

**文档版本：** v1.0.0
**最后更新：** 2026-03-08
**维护者：** Jesse Duffield 及全体贡献者
