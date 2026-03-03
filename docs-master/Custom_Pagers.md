# 自定义分页器

Lazygit 支持自定义分页器，在 config.yml 配置文件中[配置](/docs/Config.md)（可以通过在状态面板中按 `e` 打开）。

Windows 用户不支持此功能，因为我们使用的包不支持 Windows。但是，请参阅[下文](#在-windows-上模拟自定义分页器)了解解决方法。

支持多个分页器；你可以使用 `|` 键在它们之间循环。如果你通常更喜欢特定的分页器，但想对某些类型的差异使用不同的分页器，这会很有用。

分页器在 git 部分的 `pagers` 数组中配置；这是一个多分页器设置的示例：

```yaml
git:
  pagers:
    - pager: delta --dark --paging=never
    - pager: ydiff -p cat -s --wrap --width={{columnWidth}}
      colorArg: never
    - externalDiffCommand: difft --color=always
```

`colorArg` 键用于设置你的 `git diff` 命令中是否需要 `--color=always` 参数。有些分页器希望它设置为 `always`，其他分页器希望它设置为 `never`。默认值是 `always`，因为这是大多数分页器需要的。

## Delta:

```yaml
git:
  pagers:
    - pager: delta --dark --paging=never
```

![](https://i.imgur.com/QJpQkF3.png)

delta 的一个很酷的功能是 --hyperlinks，它为左边距的行号渲染可点击的链接，lazygit 支持这些。要使用它们，将 `pager:` 配置设置为 `delta --dark --paging=never --line-numbers --hyperlinks --hyperlinks-file-link-format="lazygit-edit://{path}:{line}"`；这允许你点击差异中带下划线的行号直接跳转到编辑器中的同一行。

请注意，由于技术原因，delta 的 `--navigate` 选项在 lazygit 中不起作用。

## Diff-so-fancy

```yaml
git:
  pagers:
    - pager: diff-so-fancy
```

![](https://i.imgur.com/rjH1TpT.png)

## ydiff

```yaml
gui:
  sidePanelWidth: 0.2 # 为你提供更多空间来并排显示内容
git:
  pagers:
    - colorArg: never
      pager: ydiff -p cat -s --wrap --width={{columnWidth}}
```

![](https://i.imgur.com/vaa8z0H.png)

小心这个，我认为 homebrew 和 pip 版本落后于 master。我需要直接下载 ydiff 脚本才能使无分页器功能正常工作。

## 使用外部差异命令

一些差异工具不能像上面那些简单的分页器那样工作，因为它们需要访问整个差异，所以仅仅后处理 git 的差异对它们来说是不够的。最著名的例子可能是 [difftastic](https://difftastic.wilfred.me.uk)。

这些可以通过使用 `externalDiffCommand` 配置在 lazygit 中使用；在 difftastic 的情况下，可以是

```yaml
git:
  pagers:
    - externalDiffCommand: difft --color=always
```

在这种情况下不使用 `colorArg` 和 `pager` 选项。

你可以为你的差异工具添加任何你喜欢的额外参数；例如

```yaml
git:
  pagers:
    - externalDiffCommand: difft --color=always --display=inline --syntax-highlight=off
```

除了在 lazygit 的 `externalDiffCommand` 配置中设置此命令外，你还可以告诉 lazygit 使用在 git 本身（`diff.external`）中配置的外部差异命令，方法是使用

```yaml
git:
  pagers:
    - useExternalDiffGitConfig: true
```

如果你还想在命令行上使用它进行差异比较，这会很有用，它还有一个优点，你可以在 `.gitattributes` 中按文件类型配置它；请参阅 https://git-scm.com/docs/gitattributes#_defining_an_external_diff_driver。

## 在 Windows 上模拟自定义分页器

有一个技巧可以使用配置为外部差异命令的 Powershell 脚本在 Windows 上模拟自定义分页器。它不完美，但肯定比什么都没有好。要做到这一点，将以下脚本保存为 `lazygit-pager.ps1` 在磁盘上的方便位置：

```pwsh
#!/usr/bin/env pwsh

$old = $args[1].Replace('\', '/')
$new = $args[4].Replace('\', '/')
$path = $args[0]
git diff --no-index --no-ext-diff $old $new
  | %{ $_.Replace($old, $path).Replace($new, $path) }
  | delta --width=$env:LAZYGIT_COLUMNS
```

在脚本的最后一行使用你选择的分页器和你喜欢的参数。就我个人而言，如果没有 delta 的 `--hyperlinks --hyperlinks-file-link-format="lazygit-edit://{path}:{line}"` 参数，我不想再使用 lazygit 了，请参阅[上文](#delta)。

在你的 lazygit 配置中，使用

```yml
git:
  pagers:
    - externalDiffCommand: "C:/wherever/lazygit-pager.ps1"
```

与"真正的"分页器相比，这种方法的主要限制是重命名不能正确显示；它们显示为好像是对旧文件的修改。（这仅影响块头；差异本身始终是正确的。）
