/**
 * @Author: Jason
 * @Description:
 * @File: system.go
 * @Date: 2022/9/27 10:58
 **/

package common

import (
	"fmt"
	"github.com/dyowoo/easysocket"
	"github.com/dyowoo/easysocket/example/common/ProtoMsg"
	"google.golang.org/protobuf/proto"
)

func ServiceSendMsg(request easysocket.IRequest, msgId int32, message proto.Message) {
	data, err := proto.Marshal(message)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var connId = uint32(0)

	property, _ := request.GetSession().GetProperty("ConnID")

	if property != nil {
		connId = property.(uint32)
	}

	dataTransfer := &ProtoMsg.S2G_DataTransfer{
		MsgID:  msgId,
		ConnId: connId,
		Data:   data,
	}

	_ = request.SendMsg(int32(ProtoMsg.CMD_SERVICE_S2G_DATA_TRANSFER), dataTransfer)
}
