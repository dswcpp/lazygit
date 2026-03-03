# 了解 Lazygit 何时忙碌/空闲

## 用例

这个主题值得有自己的文档，因为它有几个接触点。我们有一个用例需要知道 Lazygit 何时空闲或忙碌，因为集成测试遵循以下过程：
1) 按下一个键
2) 等待直到 Lazygit 空闲
3) 运行断言 / 按下另一个键
4) 重复

过去的过程是：
1) 按下一个键
2) 运行断言
3) 如果断言失败，等待一会儿然后重试
4) 重复

旧过程存在问题，因为由于某些视图的内容自上次按键以来尚未更改，断言可能会给出误报。

## 解决方案

首先，重要的是要区分三种不同类型的 goroutine：
* UI goroutine，只有一个，它无限地处理事件队列
* Worker goroutine，它们执行一些工作，然后通常在 UI goroutine 中排队一个事件以显示结果
* Background goroutine，它们定期生成 worker goroutine（例如每分钟执行一次 git fetch）

区分 worker goroutine 和 background goroutine 的意义在于，当任何 worker goroutine 运行时，我们认为 Lazygit 是"忙碌"的，而 background goroutine 则不是这种情况。让 background goroutine 被视为"忙碌"是没有意义的，因为那样 Lazygit 在整个程序期间都会被视为忙碌！

在 gocui 中，我们用于管理 UI 和事件的底层包，我们使用 `Task` 类型跟踪有多少忙碌的 goroutine。任务表示 lazygit 正在执行的一些工作。gocui Gui 结构体持有一个任务映射，并允许创建新任务（将其添加到映射中）、暂停/继续任务以及将任务标记为完成（将其从映射中删除）。只要映射中至少有一个忙碌的任务，Lazygit 就被认为是忙碌的；否则它被认为是空闲的。当 Lazygit 从忙碌变为空闲时，它会通知集成测试。

重要的是我们遵守以下规则，以确保在用户执行任何操作后，所有后续处理都在一个连续的忙碌块中进行，没有间隙。

### 生成 worker goroutine

这是 `OnWorker` 的基本实现（使用与 `WaitGroup` 相同的流程）：

```go
func (g *Gui) OnWorker(f func(*Task)) {
	task := g.NewTask()
	go func() {
		f(task)
		task.Done()
	}()
}
```

这里的关键是我们在生成 goroutine _之前_创建任务，因为这意味着在 goroutine 完成之前，我们在映射中至少有一个忙碌的任务。如果我们在 goroutine 内创建任务，当前函数可能会退出，Lazygit 会在 goroutine 启动之前被视为空闲，导致我们的集成测试过早进行。

你通常使用 `self.c.OnWorker(f)` 调用它。请注意，回调函数接收任务。这允许回调暂停/继续任务（见下文）。

### 生成 background goroutine

生成 background goroutine 很简单：

```go
go utils.Safe(f)
```

其中 `utils.Safe` 是一个辅助函数，确保如果 goroutine 发生 panic，我们会清理 gui。

### 以编程方式排队 UI 事件

这通过 `self.c.OnUIThread(f)` 调用。在内部，它在将函数作为事件排队之前创建一个任务（在事件结构体中包含任务），一旦该事件被事件队列处理（以及任何其他待处理的事件被处理），任务就会通过调用 `task.Done()` 从映射中删除。

### 按键

如果用户按下一个键，事件将自动排队，并在事件处理之前（和之后 `Done`）创建一个任务。

## 特殊情况

有几个特殊情况，我们直接在客户端代码中手动暂停/继续任务。这些可能会更改，但为了完整性：

### 写入主视图

如果用户在文件面板中聚焦一个文件，我们为该文件运行 `git diff` 命令并将输出写入主视图。但我们只读取足够填充视图视口的命令输出：进一步加载仅在用户滚动时发生。鉴于我们有一个后台 goroutine 用于运行命令并在滚动时写入更多输出，我们创建自己的任务并在视口填充后立即对其调用 `Done`。

### 从 git 命令请求凭据

某些 git 命令（例如 git push）可能会请求凭据。这与上面的情况相同；我们使用 worker goroutine 并在从等待 git 命令到等待用户输入时手动暂停继续其任务。这需要将任务传递给 `Push` 方法，以便可以暂停/继续它。
