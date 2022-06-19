# Echo cli
echo-cli可以让你以更加可控的方式来使用，在整个过程中你会了解到quic的一些特性。在这里我也会讨论一些在编程中遇到的细节问题。
同时在客户端和服务器端运行`go run main.go`语句即可运行quic的示例程序。
## 密钥配对与验证
`echo`服务器端会自动生成`tls.Config`文件，在这一节中我们会对其进行简要介绍。在使用TLS证书的过程中会生成成对的公钥和私钥，它使用了一种非对称式的加密算法来对传输的数据进行加密。
1.服务器端将公钥发送给客户端，CA对公钥进行验证
2.客户端生成密钥来对数据传输进行加密，之后将其发送给服务器端。
3.服务器端对密钥解密拿到客户端生成的密钥。
## 对io.copy进行封装
服务器端实现了`loggingWriter`来对`io.writer`进行封装，并将其用于打印信息。例如如果我们想要在向流写入数据之前先打印信息，我们可以使用以下代码：
```go
var stream quic.Stream
data:="Data I want to write"
fmt.Printf("Write data: %s to stream\n",data)
stream.Write([]byte(data))
```
上面的代码实现了打印信息的目的，然而其具有以下几个缺点：
- 向quic流中写入和向标椎输出流写入二者无关，这两者如果存在联系的话会更加易于管理。
- 当项目过大的时候会难以维护代码。两者都是向流中写入数据但是却接口不同，因此很容易出现漏掉某一个写入操作的情况。
因此，更好的方式是对其进行封装，使用同一接口来对不同的流执行写入操作：
```go
type loggingWriter struct {
	io.Writer
}

func (w loggingWriter) Write(b []byte)  (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}
```
## 如果服务器端使用接连两次的写操作来发送信息会怎么样呢？
使用上述的封装函数：
```go
func (w loggingWriter) Write(b []byte)  (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	if w.writeType=="twice"{
		l:=len(b)
		n,err:=w.Writer.Write(b[:l/2])
		nn,err:=w.Writer.Write(b[l/2:])
		return n+nn,err
	}
	return w.Writer.Write(b)
}
```
经过测试，并不影响客户端正常接收信息。
## OpenStream和OpenStreamSync之间的区别
从文档中我们可以看到:
- OpenStream
  OpenStream打开一个双向的QUIC流。这一过程并不会告知与其连接的另外一端：当数据流已经发送的时候，另一端只能够接收。
  如果`error`被置位，它满足了`net.Error`接口的要求。
  当达到另一端数据流的最大限度时，`err.Temperary`将会返回True。
  如果连接因为超时而关闭，`Timeout()`将会返回True。
- OpenStreamSync
  OpenStreamSync同样打开一个双向的QUIC流。在一个新的流被打开之前它会处于阻塞状态。如果`error`被置位，它满足了`net.Error`接口的要求。
  如果连接因为超时而关闭，`Timeout()`将会返回True。
两者的区别在于另一端是否知道并且接受新创建的数据流。在echo场景下，两者没有进一步的区别。
## 当服务器端结束了对流的写入会发生什么呢？
参考RFC9000，双向流的状态图如下所示：
```go
o
| Create Stream (Sending)
| Peer Creates Bidirectional Stream
v
+-------+
| Ready | Send RESET_STREAM
|       |-----------------------.
+-------+                       |
|                           |
| Send STREAM /             |
|      STREAM_DATA_BLOCKED  |
v                           |
+-------+                       |
| Send  | Send RESET_STREAM     |
|       |---------------------->|
+-------+                       |
|                           |
| Send STREAM + FIN         |
v                           v
+-------+                   +-------+
| Data  | Send RESET_STREAM | Reset |
| Sent  |------------------>| Sent  |
+-------+                   +-------+
|                           |
| Recv All ACKs             | Recv ACK
v                           v
+-------+                   +-------+
| Data  |                   | Reset |
| Recvd |                   | Recvd |
+-------+                   +-------+
```
RFC9000中是这样进行描述的：
> 在"send(发送中)"状态下, 连接端会以STEAM帧的形式对数据流进行传输(如果需要的话就重传).连接端始终将速度维持在另一端所设定的流量控制的极限值, 接受并处理`MAX_STREAM_DATA`大小的数据帧. 如果超出了数据流发送的极限值连接端会生成`STREAM_DATA_BLOCKED`数据帧, 从"send(发送中)"状态转化为阻塞状态.
> 在应用确定所有的数据都已方发送后, 另一个包含着FIN的STREAM数据帧将会发送过去, 发送端进入"Data sent"状态.

因此在echo的场景下, 应用如果不表明"所有数据都已被发送"的话, 服务器端会一直处于"send(发送中)"状态.
所以服务器必须等待定时器超时(io.copy一直等待流中数据直到遇到EOF), 之后才能进入"sent(发送完毕)"状态.
## 为什么服务器只能处理一个客户端的请求？
第一个客户的请求得到了echo的信息，但第二个请求进入阻塞直到出现超时错误:
no recent network activity "。 
那么，为什么服务器没有退出，却不能为一个新的客户请求提供服务呢？

原因是我们的服务器只为第一个流设置逻辑。如果客户端运行两次，第二次就会尝试创建一个新的流。然而，在服务器端没有对应的业务逻辑来解决这个问题，在这里的echo例子中，请求失败。
