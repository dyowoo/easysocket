/**
 * @Author: Jason
 * @Description:
 * @File: message
 * @Version: 1.0.0
 * @Date: 2022/4/7 9:39
 */

package easysocket

type IMessage interface {
	GetMsgId() int32
	GetDataLen() uint32
	GetData() []byte
	SetMsgId(msgType int32)
	SetDataLen(len uint32)
	SetData(data []byte)
}

type Message struct {
	id      int32
	dataLen uint32
	data    []byte
}

func NewMessage(id int32, data []byte) *Message {
	return &Message{
		id:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
}

func (m *Message) GetMsgId() int32 {
	return m.id
}

func (m *Message) GetDataLen() uint32 {
	return m.dataLen
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetMsgId(id int32) {
	m.id = id
}

func (m *Message) SetDataLen(len uint32) {
	m.dataLen = len
}

func (m *Message) SetData(data []byte) {
	m.data = data
}
