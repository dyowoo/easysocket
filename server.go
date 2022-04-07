/**
 * @Author: Jason
 * @Description:
 * @File: server
 * @Version: 1.0.0
 * @Date: 2022/4/7 11:14
 */

package easysocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

type HookFunc func(session ISession)

type IServer interface {
	start()
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

var connId uint32 = 0

type Server struct {
	serverName  string
	ServerType  ServerType
	host        string
	port        int
	msgHandle   IMessageHandler
	sessMgr     ISessionManager
	OnConnStart HookFunc
	OnConnStop  HookFunc
	dataPack    IDataPack
}

func NewServer(name string, serverType ServerType, host string, port int) *Server {
	return &Server{
		serverName: name,
		ServerType: serverType,
		host:       host,
		port:       port,
		msgHandle:  NewMessageHandler(),
		sessMgr:    NewSessionManager(),
		dataPack:   NewDataPack(),
	}
}

func (s *Server) startTCPServer() {
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
		if err != nil {
			fmt.Println("listen tcp error:", err.Error())
			return
		}

		fmt.Println("start ", s.serverName, " success...")

		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("accept error:", err.Error())
				continue
			}

			fmt.Println("new client connect, remote addr = ", conn.RemoteAddr().String())

			connId++
			sess := NewTCPSession(s, conn, connId, s.msgHandle)

			go sess.Start()
		}
	}()
}

func (s *Server) startWsServer() {
	go func() {
		upgrade := &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			conn, err := upgrade.Upgrade(writer, request, nil)

			if err != nil {
				fmt.Println("websocket error:", err)
				return
			}

			connId++
			sess := NewWsSession(s, conn, connId, s.msgHandle)

			s.sessMgr.Add(sess)

			go sess.Start()
		})

		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), nil); err != nil {
			fmt.Println("http listen error:", err)
			return
		}

		fmt.Println("websocket server is running...")
	}()
}

func (s *Server) startTCPClient() {
	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
		if err != nil {
			fmt.Println("connect error", err)
			return
		}
		connId++
		sess := NewTCPSession(s, conn, connId, s.msgHandle)
		go sess.Start()
	}()
}

func (s *Server) start() {
	go s.msgHandle.StartWorkerPool()

	switch s.ServerType {
	case TcpServer:
		s.startTCPServer()
	case WsServer:
		s.startWsServer()
	case TcpClient:
		s.startTCPClient()
	}

	fmt.Printf("%s at %s:%d 已启动...\n", s.serverName, s.host, s.port)
}

func (s *Server) Stop() {
	s.sessMgr.Clear()
}

func (s *Server) Serve() {
	s.start()
}

func (s *Server) AddRouter(msgId int32, router IRouter) {
	s.msgHandle.AddRouter(msgId, router)
}

func (s *Server) GetSessMgr() ISessionManager {
	return s.sessMgr
}

func (s *Server) SetOnConnStart(hookFunc HookFunc) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc HookFunc) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(session ISession) {
	if s.OnConnStart != nil {
		s.OnConnStart(session)
	}
}

func (s *Server) CallOnConnStop(session ISession) {
	if s.OnConnStop != nil {
		s.OnConnStop(session)
	}
}

func (s *Server) DataPack() IDataPack {
	return s.dataPack
}
