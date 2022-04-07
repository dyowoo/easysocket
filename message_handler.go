/**
 * @Author: Jason
 * @Description:
 * @File: message_handle
 * @Version: 1.0.0
 * @Date: 2022/4/7 14:41
 */

package easysocket

import (
	"fmt"
	"strconv"
)

type IMessageHandler interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgId int32, router IRouter)
	startOneWorker(workerId int, taskQueue chan IRequest)
	SendMsgToTaskQueue(request IRequest)
}

type MessageHandler struct {
	routers        map[int32]IRouter
	workerPoolSize uint32
	taskQueue      []chan IRequest
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		routers:        make(map[int32]IRouter),
		workerPoolSize: 10,
		taskQueue:      make([]chan IRequest, 10),
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由worker进行处理
func (m *MessageHandler) SendMsgToTaskQueue(request IRequest) {
	workerId := request.GetSession().GetConnId() % m.workerPoolSize
	m.taskQueue[workerId] <- request
}

// DoMsgHandler 处理消息
func (m *MessageHandler) DoMsgHandler(request IRequest) {
	handler, ok := m.routers[request.GetMsgId()]

	if !ok {
		fmt.Println("router msgId = ", request.GetMsgId(), " is not found")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 添加具体消息处理逻辑
func (m *MessageHandler) AddRouter(msgId int32, router IRouter) {
	if _, ok := m.routers[msgId]; ok {
		panic("repeated router, msgId = " + strconv.Itoa(int(msgId)))
	}

	m.routers[msgId] = router
}

// 启动一个worker工作进程
func (m *MessageHandler) startOneWorker(workerId int, taskQueue chan IRequest) {
	fmt.Println("worker ID = ", workerId, " is started.")
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// StartWorkerPool 启动工作池
func (m *MessageHandler) StartWorkerPool() {
	for i := 0; i < int(m.workerPoolSize); i++ {
		m.taskQueue[i] = make(chan IRequest, 1024)
		go m.startOneWorker(i, m.taskQueue[i])
	}
}
