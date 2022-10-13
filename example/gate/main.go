/**
 * @Author: Jason
 * @Description:
 * @File: main.go
 * @Date: 2022/9/26 15:32
 **/

package main

import (
	"fmt"
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/common/ProtoMsg"
	"github.com/dyowoo/easysocket/example/gate/managers"
	"github.com/dyowoo/easysocket/example/gate/routers"
	"google.golang.org/protobuf/proto"
)

var CS easysocket.IServer
var GS easysocket.IServer

func main() {
	// 处理客户端
	CS = easysocket.NewServer("ClientGate", easysocket.TcpServer, "0.0.0.0", 19000)

	//CS.SetGateHandler(gateHandler)
	CS.AddRouter(int32(ProtoMsg.CMD_GAME_PING), &routers.PingRouter{}, ProtoMsg.C2S_Ping{})
	CS.SetOnConnStart(playerStart)
	CS.SetOnConnStop(playerStop)

	CS.Serve()

	// 处理服务器
	GS = easysocket.NewServer("ServerGate", easysocket.TcpServer, "0.0.0.0", 19001)
	GS.AddRouter(int32(ProtoMsg.CMD_SERVICE_S2G_DATA_TRANSFER), &routers.DataTransferRouter{}, ProtoMsg.S2G_DataTransfer{})

	GS.SetOnConnStart(serverRegister)
	GS.SetOnConnStop(serverStop)

	GS.Serve()

	select {}

}

func playerStart(session easysocket.ISession) {
	fmt.Println("新玩家连接: ", session.RemoteIP())
	managers.PlayerMgr.Add(&managers.Player{
		Session: session,
	})
}

func playerStop(session easysocket.ISession) {
	fmt.Println("玩家断开连接")
	managers.PlayerMgr.Remove(session)
}

func serverRegister(session easysocket.ISession) {
	fmt.Println("注册服务")
	serverItem := &managers.ServerItem{
		Session: session,
		Count:   0,
	}

	managers.ServerMgr.Add(serverItem)
}

func serverStop(session easysocket.ISession) {
	fmt.Println("服务器掉线")

	managers.ServerMgr.Remove(session)
}

func gateHandler(request easysocket.IRequest) {
	dataTransfer := &ProtoMsg.G2S_DataTransfer{
		ConnId: request.GetSession().GetConnId(),
		MsgID:  request.GetMsgId(),
		Data:   request.GetData(),
	}

	var serverID = uint32(0)

	property, _ := request.GetSession().GetProperty("ServerID")

	if property != nil {
		serverID = property.(uint32)
	}

	serverItem := managers.ServerMgr.Get(serverID)

	if serverItem == nil {
		fmt.Println("没有服务器在线...")
		return
	}

	if property == nil {
		_ = request.GetSession().SetProperty("ServerID", serverItem.Session.GetConnId())
	}

	buffer, _ := proto.Marshal(dataTransfer)

	err := serverItem.Session.SendBuffMsg(int32(ProtoMsg.CMD_SERVICE_G2S_DATA_TRANSFER), buffer)
	if err != nil {
		fmt.Println(err.Error())
	}
}
