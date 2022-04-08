/**
 * @Author: Jason
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2022/4/8 9:56
 */

package main

import (
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/server/ProtoMsg"
	"github.com/dyowoo/easysocket/example/server/routers"
)

func main() {
	s := easysocket.NewServer("server", easysocket.TcpServer, "127.0.0.1", 29000)

	s.AddRouter(int32(ProtoMsg.CMD_PING), &routers.PingRouter{}, ProtoMsg.C2S_Ping{})

	go s.Serve()

	select {}
}
