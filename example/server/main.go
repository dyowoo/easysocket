/**
 * @Author: Jason
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2022/4/8 9:56
 */

package main

import (
	"fmt"
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/common/ProtoMsg"
	"github.com/dyowoo/easysocket/example/server/routers"
	"google.golang.org/protobuf/proto"
	"time"
)

var server easysocket.IServer

func main() {
	server = easysocket.NewServer("G2SServer", easysocket.TcpClient, "127.0.0.1", 19001)

	server.SetOnConnStop(OnConnStop)

	server.AddRouter(int32(ProtoMsg.CMD_SERVICE_G2S_DATA_TRANSFER), &DataTransferRouter{}, ProtoMsg.G2S_DataTransfer{})
	server.AddRouter(int32(ProtoMsg.CMD_GAME_PING), &routers.PingRouter{}, ProtoMsg.C2S_Ping{})
	server.Serve()

	fmt.Println("连接网关服务器")

	select {}
}

func OnConnStop(session easysocket.ISession) {
	fmt.Println("与网关断开连接, 3秒后重新连接")
	time.Sleep(3 * time.Second)
	server.Serve()
}

type DataTransferRouter struct {
	easysocket.BaseRouter
}

func (r *DataTransferRouter) Handle(request easysocket.IRequest, message proto.Message) {
	msg := message.(*ProtoMsg.G2S_DataTransfer)
	_ = request.GetSession().SetProperty("ConnID", msg.ConnId)

	request.SetMsg(msg.GetMsgID(), msg.GetData())

	server.SendBufferMsg(request)
}
