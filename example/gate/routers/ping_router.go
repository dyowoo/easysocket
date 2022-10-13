package routers

import (
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/common"
	"github.com/dyowoo/easysocket/example/common/ProtoMsg"
	"google.golang.org/protobuf/proto"
)

type PingRouter struct {
	easysocket.BaseRouter
}

func (r *PingRouter) Handle(request easysocket.IRequest, message proto.Message) {
	msg := message.(*ProtoMsg.C2S_Ping)

	pong := &ProtoMsg.S2C_Pong{
		Pong: msg.Ping,
	}

	common.ServiceSendMsg(request, int32(ProtoMsg.CMD_GAME_PONG), pong)
}
