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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"reflect"
)

type PreRouterHandle func(request IRequest, message proto.Message) bool

type IMessageHandler interface {
	DoMsgHandler(request IRequest)
	SetGateHandler(handler GateHandler)
	AddPreRouter(handle PreRouterHandle)
	AddRouter(msgId int32, router IRouter, v any)
	startOneWorker(workerId int, taskQueue chan IRequest)
	SendMsgToTaskQueue(request IRequest)
	StartWorkerPool()
}

type MessageHandler struct {
	gateHandler     GateHandler
	preRouterHandle PreRouterHandle
	routers         map[int32]IRouter
	protocols       map[int32]string
	workerPoolSize  uint32
	taskQueue       []chan IRequest
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		gateHandler:     nil,
		preRouterHandle: nil,
		routers:         make(map[int32]IRouter),
		protocols:       make(map[int32]string),
		workerPoolSize:  10,
		taskQueue:       make([]chan IRequest, 10),
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由worker进行处理
func (m *MessageHandler) SendMsgToTaskQueue(request IRequest) {
	workerId := request.GetSession().GetConnId() % m.workerPoolSize
	m.taskQueue[workerId] <- request
}

// ReflectProto 通过反射把数据解析成proto message
func (m *MessageHandler) ReflectProto(request IRequest) proto.Message {
	msgId := request.GetMsgId()
	data := request.GetData()

	if _, ok := m.protocols[msgId]; !ok {
		fmt.Println("msgId: ", msgId, " not exist")
		return nil
	}

	msgRef, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(m.protocols[msgId]))

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	msg := msgRef.New().Interface().(proto.Message)

	_ = proto.Unmarshal(data, msg)

	return msg
}

// DoMsgHandler 处理消息
func (m *MessageHandler) DoMsgHandler(request IRequest) {
	if m.gateHandler != nil {
		m.gateHandler(request)
	} else {
		handler, ok := m.routers[request.GetMsgId()]

		if !ok {
			fmt.Println("router msgId = ", request.GetMsgId(), " is not found")
			return
		}

		msg := m.ReflectProto(request)

		if msg == nil {
			return
		}

		if m.preRouterHandle != nil && !m.preRouterHandle(request, msg) {
			return
		}

		handler.PreHandle(request, msg)
		handler.Handle(request, msg)
		handler.PostHandle(request, msg)
	}
}

// SetGateHandler 设置网关处理函数
func (m *MessageHandler) SetGateHandler(handler GateHandler) {
	m.gateHandler = handler
}

// AddPreRouter 添加路由前置处理
func (m *MessageHandler) AddPreRouter(handle PreRouterHandle) {
	m.preRouterHandle = handle
}

// AddRouter 添加具体消息处理逻辑
func (m *MessageHandler) AddRouter(msgId int32, router IRouter, v any) {
	if _, ok := m.routers[msgId]; ok {
		panic(fmt.Sprintf("repeated router, msgId = %d", msgId))
	}

	if _, ok := m.protocols[msgId]; ok {
		panic(fmt.Sprintf("repeated protocol, msgId = %d", msgId))
	}

	m.routers[msgId] = router
	m.protocols[msgId] = reflect.TypeOf(v).Name()
}

// 启动一个worker工作进程
func (m *MessageHandler) startOneWorker(workerId int, taskQueue chan IRequest) {
	//fmt.Println("worker ID = ", workerId, " is started.")
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// StartWorkerPool 启动工作池
func (m *MessageHandler) StartWorkerPool() {
	fmt.Println("启动工作池...")
	for i := 0; i < int(m.workerPoolSize); i++ {
		m.taskQueue[i] = make(chan IRequest, 1024)
		go m.startOneWorker(i, m.taskQueue[i])
	}
}
