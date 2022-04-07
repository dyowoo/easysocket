/**
 * @Author: Jason
 * @Description:
 * @File: request
 * @Version: 1.0.0
 * @Date: 2022/4/7 11:06
 */

package easysocket

type IRequest interface {
	GetSession() ISession
	GetData() []byte
	GetMsgId() int32
}

type Request struct {
	sess ISession
	msg  IMessage
}

// GetSession 获取请求连接信息
func (r *Request) GetSession() ISession {
	return r.sess
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgId 获取请求的消息的ID
func (r *Request) GetMsgId() int32 {
	return r.msg.GetMsgId()
}
