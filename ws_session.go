/**
 * @Author: Jason
 * @Description:
 * @File: ws_session
 * @Version: 1.0.0
 * @Date: 2022/4/7 16:29
 */

package easysocket

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
)

type WsSession struct {
	Session
	conn *websocket.Conn
}

func NewWsSession(server IServer, conn *websocket.Conn, connId uint32, handler IMessageHandler) *WsSession {
	sess := &WsSession{
		Session: Session{
			server:      server,
			connId:      connId,
			msgHandle:   handler,
			msgBuffChan: make(chan []byte, 1024),
			property:    nil,
			isClosed:    false,
		},
		conn: conn,
	}

	sess.server.GetSessMgr().Add(sess)

	return sess
}

func (s *WsSession) startReader() {
	fmt.Println("Reader goroutine is running...")
	defer fmt.Println(s.RemoteAddr().String(), " conn reader exit!")
	defer s.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			_, data, err := s.conn.ReadMessage()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// 获取消息头数据
			headData := data[:DP.GetHeadLen()]
			// 解析消息头
			msg := DP.UnPack(headData)

			if msg.GetDataLen() > 0 {
				msg.SetData(data[DP.GetHeadLen():])

				req := &Request{
					sess: s,
					msg:  msg,
				}

				s.msgHandle.SendMsgToTaskQueue(req)
			}
		}
	}
}

func (s *WsSession) startWriter() {
	fmt.Println("Writer goroutine is running...")
	defer fmt.Println(s.RemoteAddr().String(), " conn writer exit!")

	for {
		select {
		case data, ok := <-s.msgBuffChan:
			if ok {
				if err := s.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
					fmt.Println("send buff data error:", err, " conn writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is closed")
				break
			}
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *WsSession) Start() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	go s.startReader()
	go s.startWriter()

	s.server.CallOnConnStart(s)

	select {
	case <-s.ctx.Done():
		s.finalizer()
		return
	}
}

func (s *WsSession) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *WsSession) SendMsg(msgId int32, data []byte) error {
	s.RLock()
	defer s.RUnlock()

	if s.isClosed {
		return errors.New("connection closed when send msg")
	}

	msg := DP.Pack(NewMessage(msgId, data))

	return s.conn.WriteMessage(websocket.BinaryMessage, msg)
}

func (s *WsSession) finalizer() {
	s.server.CallOnConnStop(s)

	s.Lock()
	defer s.Unlock()

	if s.isClosed {
		return
	}

	fmt.Println("conn stop()...connId = ", s.connId)

	s.isClosed = true

	_ = s.conn.Close()

	s.server.GetSessMgr().Remove(s)

	close(s.msgBuffChan)
}
