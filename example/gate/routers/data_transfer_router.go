/**
 * @Author: Jason
 * @Description:
 * @File: data_transfer_router.go
 * @Date: 2022/9/26 15:53
 **/

package routers

import (
	"fmt"
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/common/ProtoMsg"
	"github.com/dyowoo/easysocket/example/gate/managers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type DataTransferRouter struct {
	easysocket.BaseRouter
}

func (r *DataTransferRouter) Handle(request easysocket.IRequest, message proto.Message) {
	msg := message.(*ProtoMsg.S2G_DataTransfer)

	player := managers.PlayerMgr.Get(msg.ConnId)

	if player == nil {
		return
	}

	msgRef, err := protoregistry.GlobalTypes.FindMessageByName("S2C_Pong")

	if err != nil {
		fmt.Println(err.Error())
	}

	m := msgRef.New().Interface().(proto.Message)

	_ = proto.Unmarshal(msg.GetData(), m)

	_ = player.Session.SendMsg(msg.MsgID, msg.Data)

}
