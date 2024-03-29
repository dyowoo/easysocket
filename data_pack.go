/**
 * @Author: Jason
 * @Description:
 * @File: data_pack
 * @Version: 1.0.0
 * @Date: 2022/4/7 9:38
 */

package easysocket

import "encoding/binary"

// DataPack 封包拆包类实例
type DataPack struct{}

// GetHeadLen 获取包头长度
func (dp *DataPack) GetHeadLen() uint32 {
	return MsgLenSize + MsgTypeSize
}

// Pack 封包 长度|类型|内容
func (dp *DataPack) Pack(message IMessage) []byte {
	headLen := dp.GetHeadLen()
	dataLen := message.GetDataLen()

	data := make([]byte, headLen+dataLen)

	// 写入消息长度
	binary.BigEndian.PutUint32(data, dataLen)
	// 写入消息类型
	binary.BigEndian.PutUint32(data[MsgLenSize:], uint32(message.GetMsgId()))
	// 写入消息内容
	_ = copy(data[headLen:], message.GetData())

	return data
}

// UnPack 拆包 先读取消息头信息 长度|类型
func (dp *DataPack) UnPack(data []byte) IMessage {
	dataLen := binary.BigEndian.Uint32(data)
	msgId := binary.BigEndian.Uint32(data[MsgLenSize:])

	msg := &Message{}

	msg.SetDataLen(dataLen)
	msg.SetMsgId(int32(msgId))

	return msg
}

var DP = DataPack{}
