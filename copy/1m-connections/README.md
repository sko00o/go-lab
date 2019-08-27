# ref

* [A Million WebSockets and Go](https://www.freecodecamp.org/news/million-websockets-and-go-cc58418460bb/)
* [gobwas/ws](https://github.com/gobwas/ws)
* [Going Infinite, handling 1M websockets connections in Go](https://speakerdeck.com/eranyanay/going-infinite-handling-1m-websockets-connections-in-go)
* [eranyanay/1m-go-websockets](https://github.com/eranyanay/1m-go-websockets)
* [百万 Go TCP 连接的思考: epoll方式减少资源占用](https://colobu.com/2019/02/23/1m-go-tcp-connection/)
* [smallnest/1m-go-tcp-server](https://github.com/smallnest/1m-go-tcp-server)

more info

Linux I/O 多路复用有 水平触发 和  边缘触发 两种方式

水平触发：
如果文件描述符（fd）已就绪，可以非阻塞执行 I/O 操作，
此时触发通知，允许在任意时刻重复检测 I/O 状态，
select 、 poll 属于水平触发

边缘触发：
文件描述符自上次状态改变后，有新的 I/O 活动到来，此时
触发通知，在收到一个 I/O 事件通知后要尽可能多的执行 I/O
操作，如果在一次通知中没有执行完 I/O 操作，就要等到下一
次新的 I/O 活动到来才能获取到就绪的描述符。
信号驱动式 I/O 属于边缘触发

ref： https://blog.csdn.net/D_Guco/article/details/71373381

epoll 网络模型

持续有数据到来时，边缘触发效率高，要注意没有读完的情况