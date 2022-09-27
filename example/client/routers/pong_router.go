/**
 * @Author: Jason
 * @Description:
 * @File: pong_router.go
 * @Date: 2022/9/26 16:01
 **/

package routers

import (
	"fmt"
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/common/ProtoMsg"
	"google.golang.org/protobuf/proto"
)

type PongRouter struct {
	easysocket.BaseRouter
}

func (r *PongRouter) Handle(request easysocket.IRequest, message proto.Message) {
	msg := message.(*ProtoMsg.S2C_Pong)

	fmt.Println("====> Pong:", request.GetMsgId(), ", data = ", msg.GetPong())
}
