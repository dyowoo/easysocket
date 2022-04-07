/**
 * @Author: Jason
 * @Description:
 * @File: server
 * @Version: 1.0.0
 * @Date: 2022/4/7 11:14
 */

package easysocket

type HookFunc func(session ISession)

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgId int32, router IRouter)
	GetSessMgr() ISessionManager
	SetOnConnStart(hookFunc HookFunc)
	SetOnConnStop(hookFunc HookFunc)
	CallOnConnStart(session ISession)
	CallOnConnStop(session ISession)
	DataPack() IDataPack
}

type Server struct {
	serverName     string
	ServerType     ServerType
	host           string
	port           int
	messageHandler IMessageHandler
	sessMgr        ISessionManager
	OnConnStart    HookFunc
	OnConnStop     HookFunc
	dataPack       IDataPack
}

func NewServer(name string, serverType ServerType, host string, port int) *Server {
	return &Server{
		serverName:     name,
		ServerType:     serverType,
		host:           host,
		port:           port,
		messageHandler: NewMessageHandler(),
		sessMgr:        NewSessionManager(),
		dataPack:       NewDataPack(),
	}
}
