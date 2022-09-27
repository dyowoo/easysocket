/**
 * @Author: Jason
 * @Description:
 * @File: ping_router.go
 * @Date: 2022/9/27 10:47
 **/

package routers

import (
	"fmt"
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

	connId, _ := request.GetSession().GetProperty("ConnID")

	fmt.Println("ConnID = ", connId)

	pong := &ProtoMsg.S2C_Pong{
		Pong: msg.Ping,
	}

	common.ServiceSendMsg(request, int32(ProtoMsg.CMD_GAME_PONG), pong)
}
