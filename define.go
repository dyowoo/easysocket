/**
 * @Author: Jason
 * @Description:
 * @File: define
 * @Version: 1.0.0
 * @Date: 2022/4/7 10:16
 */

package easysocket

type HookFunc func(session ISession)

type ServerType uint32
type GateHandler func(request IRequest)

const (
	NullServer ServerType = iota
	TcpServer
	WsServer
	TcpClient
)

const (
	MsgTypeSize = 4
	MsgLenSize  = 4
)
