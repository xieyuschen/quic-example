# 3.`echo`运行期间的包
本节旨在帮助您了解 echo 示例中的包.
如果我们想用wireshark抓包分析，密钥对和证书应该保持不变，因为我们需要它来解密包。 因此，证书文件存储在 `/cert` 文件夹下.
## 如何对包进行解密
于2021年5月发布的QUIC标准文档中说明了wireshark可以自动对包进行分析.
## 请注意需要使用支持RFC Draft-27 的版本
当您捕获 quic 数据包时, 您可能会在 wireshark 中看到此错误信息. 这是因为您没有 TLS 的主密钥.
```
Expert Info (Warning/Decryption):
Failed to create decryption context: Secrets are not available
```
您可以配置TLS获取密钥，然后将此密钥加载到wireshark，这样你就可以看到解密的包内容.
[Issue 5](https://github.com/xieyuschen/quic-example/issues/5) 中已经记录了该问题. 
## TLS握手
[TLS reference in cloudflare blog](https://www.cloudflare.com/learning/ssl/what-happens-in-a-tls-handshake/).
TLS 是一种加密协议，旨在保护 Internet 通信. TLS 握手是启动使用 TLS 加密通信会话的过程. 在 TLS 握手期间，通信双方交换
消息相互确认，相互验证，确定他们将使用的加密算法，并确认将要使用的会话密钥。

TSL（传输层安全）是在ssl（安全套接字层）基础上的产物. 握手中包括版本协商、确保密码、交换证书和生成会话等一系列过程:

- 指定将使用的 TLS 版本（TLS 1.0、1.2、1.3 等
- 决定他们将使用哪些密码套件（见下文）
- 通过服务器的公钥和 SSL 证书机构的电子签名确定服务器身份
- 生成会话密钥以便在握手完成后使用对称加密
  
**注意，QUIC使用了基于短暂Diffie-Hellman算法的TLSv1.3版本**

### TLS握手的步骤

1. `client hello` 消息
    `client hello` 发送：
    - 客户端支持的 TLS 版本
    - 支持的密码套件
    - 一串随机字节，称为`client random`

2. `server hello` 消息
    `server hello` 回复 `client hello`：
    - 服务器的 SSL 证书
    - 服务器选择的密码套件
    - `server random`

3. 服务器的数字签名
服务器使用其私钥加密`client random`、`server random`及其 DH 参数*等.
该加密数据充当服务器的数字签名，确定服务器的拥有私钥与 SSL 证书中的公钥匹配.

4. 数字签名确认
客户端使用公钥解密服务器的数字签名，确定拥有私钥的服务器的身份.
客户端 DH 参数：客户端将其 DH 参数发送给服务器. 

5. 客户端和服务器计算预主密钥
客户端和服务器使用它们交换的 DH 参数分别计算匹配的预主密钥. 

6. 创建会话密钥
    这一步生成会话. 客户端和服务器都从`client random`, 和`server random`随机生成会话密钥，和预主密钥。 **他们应该得到相同的结果**

7. 客户准备就绪
    客户端发送使用会话密钥加密的“完成”消息

8. 服务器准备就绪
    服务器发送使用会话密钥加密的“完成”消息

9. 实现安全对称加密
    握手完成，使用会话密钥继续通信

*DH 参数：DH 代表 Diffie-Hellman. Diffie-Hellman 算法使用指数计算得出
同一个预主密钥. 服务端和客户端各提供一个参数进行计算，尽管两端进行不同的计算但是其结果相同. 

## QUIC包格式

Quic 数据包格式由 RFC 9000 标准定义，您可以查看它来查找更多信息. 这里的主题旨在简要介绍以下包分析. 

### QUIC的报头类型: 长和短的报头

与包头格式固定的 TCP 不同，QUIC有两种类型的报头。

- 长报头：
   建立连接的QUIC包需要包含若干条信息，它使用长报头格式. 
- 短报头：
   握手完成后使用的短报头格式. 一旦建立连接，只需要某些报头字段，随后的包使用短报头以提高效率. 

### 包格式

- 初始包：
- 重试包：
- 握手包：
- 受保护的有效载荷包：
  
## QUIC握手过程
在这里的 echo 示例中，客户端和服务器之间的工作流如下所示，工作流图中的原点将连接建立和流建立分开. 

![](img/echo-workflow.png)

### 第一步: 客户端发送的初始包

[RFC 17.2.2](https://www.rfc-editor.org/rfc/rfc9000.html#section-17.2.2)显示初始数据包格式. 当我们使用wireshark分析数据包，已经解组了数据包，这里就不过多解释格式了只注意以下几点:

- Quic 是基于 UDP 和 IP 层的，所以它包含 IP 包和 UDP 包。 IP 头包含 20 个字节（或 5 个32 位增量，最大 24 字节）。 UDP 标头长度为 8 字节. 注意，在十六进制显示模式下，每个十六进制字符代表 4 位（2^4）。

在第一步中，初始包包如下图所示：
![img.png](img/1-initial.png)
- 随机目标连接ID和空源连接ID：
   当客户端想要建立一个新的连接时，`Destination Connection ID`是有效的，但是`Source Connection ID` 作为 null 无效. 
   [7.2.3-Negotiating Connection IDs](https://www.rfc-editor.org/rfc/rfc9000.html#section-7.2-3)中声称:

   > 当一个初始包由一个以前没有收到过初始包或重试包的客户端发送时
   > 服务器，客户端使用不可预测的值填充目标连接 ID 字段. 这个目的地
   > 连接 ID 的长度必须至少为 8 个字节. 在从服务器接收到数据包之前，客户端必须在此连接中对所有包使用相同的目标连接 ID 值. 

    因此，在第一步中，目标连接 ID 由客户端生成. 

- 空令牌：
令牌用于验证地址以避免攻击，它是由服务器生成的.
所以在第一步中需要一个空令牌。

- CRYPTO：`client hello`
   客户端发送的第一个包总是包含一个 CRYPTO 帧，该帧包含第一个数据包的开始或全部的加密握手消息. 发送的第一个 CRYPTO 帧的偏移量总是从 0 开始.
   CRYPTO 帧直接封装了一个 TLSv1.3 `ClientHello`.

### 第二步: 服务器端发送的重试包
发送重试数据包目的在于：

- 开始确定连接ID.
- 将令牌传递给客户端.

![img.png](img/2-retry.png)

[RFC9000-7-Negotiating Connection IDs](https://www.rfc-editor.org/rfc/rfc9000.html#name-negotiating-connection-ids)中写道:

> 客户端发送的第一个初始数据包中的目标连接 ID 字段用于确定包
> 保护密钥. 这些密钥在收到重试数据包后会发生变化.

**此步骤旨在验证连接是否有效以避免流量放大攻击**并且令牌将通过“RETRY”包发送给客户端. 

#### QUIC的代码实现细节
代码实现的细节在`server.go:395 func (s *baseServer) handleInitialImpl(p *receivedPacket, hdr *wire.Header) error`中. 我们可以检查在这种情况下服务器端将发送一个`RETRY`包. 

```go
	if !s.config.AcceptToken(p.remoteAddr, token) {
		go func() {
			defer p.buffer.Release()
			if token != nil && token.IsRetryToken {
				if err := s.maybeSendInvalidToken(p, hdr); err != nil {
					s.logger.Debugf("Error sending INVALID_TOKEN error: %s", err)
				}
				return
			}
			if err := s.sendRetry(p.remoteAddr, hdr, p.info); err != nil {
				s.logger.Debugf("Error sending Retry: %s", err)
			}
		}()
		return nil
	}
```

结果，函数调用`s.config.AcceptToken(p.remoteAddr, token)`决定是否发送`RETRY`。 添加一个断点查看详情. 服务器端通过以下代码检查客户端初始数据包，这是一个相当好的抽象接口的使用(未完待续).

```go
var defaultAcceptToken = func(clientAddr net.Addr, token *Token) bool {
	if token == nil {
		return false
	}
	validity := protocol.TokenValidity
	if token.IsRetryToken {
		validity = protocol.RetryTokenValidity
	}
	if time.Now().After(token.SentTime.Add(validity)) {
		return false
	}
	var sourceAddr string
	if udpAddr, ok := clientAddr.(*net.UDPAddr); ok {
		sourceAddr = udpAddr.IP.String()
	} else {
		sourceAddr = clientAddr.String()
	}
	return sourceAddr == token.RemoteAddr
}
```
如果需要发送 `RETRY` 包，则服务器端构造一个令牌发送回去：

```go
func (s *baseServer) sendRetry(remoteAddr net.Addr, hdr *wire.Header, info *packetInfo) error {
  // ignore some lines
  token, err := s.tokenGenerator.NewRetryToken(remoteAddr, hdr.DestConnectionID, srcConnID)
  // ignore some lines
}
```
### 第三步: 客户端发送的第二个初始化包

初始包目的在于：

- 使用从 `RETRY` 获得的以下值发送初始包：
   - 重试令牌.
   - 重试的目标连接 ID.
   - `CRYPTO` 框架中的客户端问候.

### 第四步: 从服务器端发送过来的握手包

握手包的目的在于：

- 发回一个 ACK.
- 发送源连接 ID.
- 让服务器向客户端打招呼并发送多个握手消息（包括加密扩展、证书
和证书验证. 此步骤是TLS握手中的多个步骤的混合. 

根据标准[RFC 13.2.1](https://www.rfc-editor.org/rfc/rfc9000.html#section-13.1-3)中所说的, ack保证了QUIC的可靠性.

(下一步任务: 了解更多的可靠性实现细节).

### 第四步: 服务器发送握手包后受保护的有效载荷

受保护的有效载荷数据包旨在：

- 在 TLS 建立中发送带有多个握手消息（第三步中服务器的数字签名）的`CRYPTO`.
- 发送`NEW_CONNECTION_ID`.

服务器发送握手后，有一个`protected packet`. 这个`protected packet`有两个 IETF quic 数据包, 一个是`CRYPTO`，另一个是`NEW_CONNECTION_ID`. 

- 包含`CRYPTO`帧的 IETF 数据包：
   它使用长报头，其`CRYPTO`帧存储 tls 握手协议消息（包括证书验证和完成消息).

- 包含`NEW_CONNECTION_ID`的 IETF 数据包：
   此数据包使用**短报头**并包含 3 个 `NEW_CONNECTION_ID` 部分. 

### 第五步: 客户端发送的初始包

发送初始包的目的在于：

- 向服务器发送 ACK.

在服务器发送包含 tls 握手包和包含`NEW_CONNECTION_ID`的受保护的有效载荷的握手包之后，客户端发送一个带有`Initial`包的确认.

### 第五步: 从客户端发来的握手包

握手包的目的在于:
- 将ACK发送给服务器
- 将握手`finish`信息发送给给服务器

在这一步中, **已经建立了TLS连接**!

握手包中包含了ACK以及TLS握手`finish`的数据

### 第五步: 客户端发送的两个受保护的有效载荷
- 第一个发回一个ACK
- 第二个发回一个`ENTIRE_CONNECTION_ID`
  
在上述过程结束之后, 整个建立连接的过程就会结束并且此时客户端会等待接收来自服务器的ACK.

注意这两者都使用了**短报头**

### 第六步: 来自服务器端的握手包

服务器发回一个使用 **短报头** 的 ACK. 

请注意，收到此数据包后，**两点之间正式建立连接**. 

## QUIC 流握手过程
QUIC流建立在QUIC连接的基础上, 在工作流图中，建立 quic 流从 step7 开始. 

### 第七步: 从客户端发送到服务器的有效载荷
- 使用短 QUIC 头部
- 像id, 双向流等这样的流信息
- **将数据发送给服务器端**

QUIC允许用户在不建立流的情况下发送数据.
![img.png](img/send-data-before-establishing.png)

### Step8：服务器回复客户端的有效载荷
在step8中，有三个包：
- 第一个包：发送 ACK.
- 第二个包：HANDSHAKE_DONE、NEW_TOKEN 和 CRYPTO 中的`New session ticket`. 
- 第三个包：一个 ACK 和“Hello world”作为echo的信息.  
**请注意，当前流仍然处于建立的过程**. 

### Step9：客户端发回一个ACK
在收到所有三个包后，客户端会发回一个 ACK. 最终, 客户端和服务器之间的交流到此结束. 
