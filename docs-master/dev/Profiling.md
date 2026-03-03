# 性能分析 Lazygit

如果你想调查是什么导致了 CPU 或内存使用，请使用 `-profile` 命令行标志启动 lazygit。这告诉它启动一个集成的 Web 服务器来监听性能分析请求。

## 保存性能分析数据

### CPU

当 lazygit 使用 `-profile` 标志运行时，在另一个终端窗口中运行此命令来执行 CPU 性能分析并将其保存到文件：

```sh
curl -o cpu.out http://127.0.0.1:6060/debug/pprof/profile
```

默认情况下，它分析 30 秒。要更改持续时间，使用

```sh
curl -o cpu.out 'http://127.0.0.1:6060/debug/pprof/profile?seconds=60'
```

### 内存

要保存堆性能分析（包含自启动以来分配的所有内存的信息），使用

```sh
curl -o mem.out http://127.0.0.1:6060/debug/pprof/heap
```

有时获取增量日志会很有用，即查看内存使用情况从一个时间点到另一个时间点的发展情况。为此，使用

```sh
curl -o mem.out 'http://127.0.0.1:6060/debug/pprof/heap?seconds=20'
```

这将记录现在和 20 秒后之间的内存使用差异，因此它给你 20 秒的时间在 lazygit 中执行你感兴趣测量的操作。

## 查看性能分析数据

要显示性能分析数据，你可以使用 speedscope.app 或 go 附带的 pprof 工具。我更喜欢前者，因为它有更好的 UI 并且功能更强大；但是，我见过一些情况下它由于某种原因无法加载性能分析，在这种情况下，最好有 pprof 工具作为后备。

### Speedscope.app

在浏览器中访问 https://www.speedscope.app/，并将保存的性能分析拖到浏览器窗口上。有关如何导航数据，请参阅[文档](https://github.com/jlfwong/speedscope?tab=readme-ov-file#usage)。

### Pprof 工具

要查看你保存为 `cpu.out` 的性能分析，使用

```sh
go tool pprof -http=:8080 cpu.out
```

默认情况下，这会显示图形视图，我个人觉得不是很有用。从视图菜单中选择"Flame Graph"以显示数据的更有用的表示。
