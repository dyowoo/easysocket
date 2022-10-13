/**
 * @Author: Jason
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2022/4/8 12:57
 */

package main

import (
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/client/routers"
	"github.com/dyowoo/easysocket/example/common/ProtoMsg"
	"time"
)

func main() {
	createClient("client", 100)
	//createClient("client2", 101)

	select {}
}

func createClient(name string, ping int32) {
	go func() {
		c := easysocket.NewServer(name, easysocket.TcpClient, "127.0.0.1", 19000)

		c.AddRouter(int32(ProtoMsg.CMD_GAME_PONG), &routers.PongRouter{}, ProtoMsg.S2C_Pong{})

		c.Serve()

		ping := &ProtoMsg.C2S_Ping{
			Ping: ping,
		}
		time.Sleep(3 * time.Second)

		for {
			time.Sleep(10)

			c.SendMsg(int32(ProtoMsg.CMD_GAME_PING), ping)
		}
	}()
}
