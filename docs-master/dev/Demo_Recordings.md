# 演示录制

我们希望我们的演示录制保持一致，并且如果我们对 Lazygit 的 UI 进行更改，可以轻松更新。幸运的是，我们有一个现有的录制系统用于集成测试，所以我们可以利用它。

你需要熟悉集成测试的编写方式：请参阅[这里](../../pkg/integration/README.md)。

## 先决条件

理想情况下，我们会通过 docker 运行整个过程，但我们还没有让它工作。所以你需要：
```
# 用于录制
npm i -g terminalizer
# 用于 gif 压缩
npm i -g gifsicle
# 用于 mp4 转换
brew install ffmpeg

# 带图标的字体
wget https://github.com/ryanoasis/nerd-fonts/releases/download/v3.0.2/DejaVuSansMono.tar.xz && \
  tar -xf DejaVuSansMono.tar.xz -C /usr/local/share/fonts && \
  rm DejaVuSansMono.tar.xz
```

## 创建演示

演示位于 `pkg/integration/tests/demo/` 中。它们类似于常规集成测试，但有 `IsDemo: true`，这有几个效果：
* UI 的底行更安静，以便我们可以渲染字幕
* Fetch/Push/Pull 有人工延迟来模拟网络请求
* 右下角的加载器不会出现

在演示中，我们不需要像在测试中那样严格地进行断言。但仍然最好有一些基本断言，这样如果我们自动化更新演示的过程，我们就会知道其中一个是否已损坏。

在编写演示时，你可以使用与集成测试相同的流程：
* 设置仓库
* 在沙盒模式下运行演示以了解需要发生什么
* 回来编写代码使其发生

### 添加字幕

最好添加字幕来解释正在执行什么任务。使用现有演示作为指南。

### 设置 assets worktree

我们将资产（包括演示录制）存储在 `assets` 分支中，这是一个与主分支不共享历史记录的分支，纯粹用于存储资产。单独存储它们意味着我们不会用大型二进制文件堵塞代码分支。

脚本和演示定义位于代码分支中，但输出位于 assets 分支中，因此要能够从演示创建视频，你需要为 assets 分支创建一个链接的 worktree，你可以使用以下命令执行此操作：

```sh
git worktree add .worktrees/assets assets
```

输出将存储在 `.worktrees/assets/demos/` 中。我们将存储三个单独的东西：
* 录制的 yaml
* 原始 gif
* 压缩的 gif 或 mp4，取决于你选择的输出（见下文）

### 录制演示

一旦你对演示满意，你可以使用以下命令录制它：
```sh
scripts/record_demo.sh [gif|mp4] <path>
# 例如
scripts/record_demo.sh gif pkg/integration/tests/demo/interactive_rebase.go
```

~~gif 格式用于 readme 的第一个视频（它的大小更大但有自动播放和循环）~~
~~mp4 格式用于其他所有内容（无循环，需要点击，但大小更小）。~~

事实证明，你不能在仓库中存储 mp4 并从 README 链接它们，所以我们现在将在所有地方使用 gif。

### 在 README/文档中包含演示

如果你遵循了上述步骤，你将在 assets worktree 中得到输出。

在该 worktree 中，暂存所有三个输出文件并针对 assets 分支提出 PR。

然后回到代码分支，在文档中，你可以像这样嵌入录制：
```md
![Nuke working tree](../assets/demo/interactive_rebase-compressed.gif)
```

这意味着我们可以更新资产而无需更新嵌入它们的文档。
