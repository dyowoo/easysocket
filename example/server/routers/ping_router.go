/**
 * @Author: Jason
 * @Description:
 * @File: login_router
 * @Version: 1.0.0
 * @Date: 2022/4/8 12:47
 */

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
