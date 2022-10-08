# easysocket

easysocket 是一个基于Golang的轻量级并发网络服务框架，支持tcp socket和websocket。内置支持protocol buffer

```shell
go get github.com/dyowoo/easysocket
```

## Example

### server

```go
package main

import (
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/server/ProtoMsg"
	"github.com/dyowoo/easysocket/example/server/routers"
)

func main() {
	s := easysocket.NewServer("server", easysocket.TcpServer, "0.0.0.0", 29000)

	s.AddRouter(int32(ProtoMsg.CMD_PING), &routers.PingRouter{}, ProtoMsg.C2S_Ping{})

	go s.Serve()

	select {}
}
```

#### router
```go
package routers

import (
	"fmt"
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/server/ProtoMsg"
	"google.golang.org/protobuf/proto"
)

type PingRouter struct {
	easysocket.BaseRouter
}

func (r *PingRouter) Handle(request easysocket.IRequest, message proto.Message) {
	msg := message.(*ProtoMsg.C2S_Ping)

	fmt.Println("===> client msgId: ", request.GetMsgId(), " msg: ", msg.GetPing())

	pong := ProtoMsg.S2C_Pong{
		Pong: "pong",
	}

	buffer, err := proto.Marshal(proto.Message(&pong))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_ = request.GetSession().SendBuffMsg(int32(ProtoMsg.CMD_PONG), buffer)
}
```

### client

```go
package main

import (
	"fmt"
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/server/ProtoMsg"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:29000")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go func() {
		for {
			dp := easysocket.NewDataPack()

			ping := &ProtoMsg.C2S_Ping{
				Ping: "ping",
			}

			buffer, _ := proto.Marshal(ping)

			msg := dp.Pack(easysocket.NewMessage(int32(ProtoMsg.CMD_PING), buffer))

			_, err := conn.Write(msg)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			headData := make([]byte, dp.GetHeadLen())
			_, err = io.ReadFull(conn, headData)

			if err != nil {
				fmt.Println(err.Error())
				break
			}

			msgHead := dp.UnPack(headData)

			if msgHead.GetDataLen() > 0 {
				msg := msgHead.(*easysocket.Message)

				data := make([]byte, msg.GetDataLen())
				_, err := io.ReadFull(conn, data)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				msg.SetData(data)

				fmt.Println("===> ping recv msgId:", msg.GetMsgId(), ", len:", msg.GetDataLen(), ", data:", string(data))
			}

			time.Sleep(1 * time.Second)
		}
	}()

	select {}
}
```

```
upstream proxy_server {
        server 127.0.0.1:19000;
}

server {
        listen 443 ssl;
        server_name xxx.xxx.com;
        location / {
                proxy_pass http://proxy_server;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header Connection "upgrade";
                proxy_connect_timeout 30d;
                proxy_send_timeout 30d;
                proxy_read_timeout 30d;
        }
        keepalive_timeout 999999999s;
        ssl_certificate cert/xxx.pem;  #需要将cert-file-name.pem替换成已上传的证书文件的名称。
        ssl_certificate_key cert/xxx.key; #需要将cert-file-name.key替换成已上传的证书密钥文件的名称。
        ssl_session_timeout 99999999m;
        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
        #表示使用的加密套件的类型。
        ssl_protocols TLSv1 TLSv1.1 TLSv1.2; #表示使用的TLS协议的类型。
        ssl_prefer_server_ciphers on;
}

```