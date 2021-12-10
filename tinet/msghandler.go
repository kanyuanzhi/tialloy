package tinet

import (
	"github.com/kanyuanzhi/tialloy/global"
	"github.com/kanyuanzhi/tialloy/tiface"
)

type MsgHandler struct {
	Apis             map[uint32]tiface.IRouter // Apis[msgID] = handler
	ServerType       string
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	TaskQueue        []chan tiface.IRequest
}

func NewMsgHandler(serverType string) tiface.IMsgHandler {
	msgHandler := &MsgHandler{
		Apis:       make(map[uint32]tiface.IRouter),
		ServerType: serverType,
	}
	switch serverType {
	case "tcp":
		msgHandler.WorkerPoolSize = global.Object.TcpWorkerPoolSize
		msgHandler.MaxWorkerTaskLen = global.Object.TcpMaxWorkerTaskLen
	case "websocket":
		msgHandler.WorkerPoolSize = global.Object.WebsocketWorkerPoolSize
		msgHandler.MaxWorkerTaskLen = global.Object.WebsocketMaxWorkerTaskLen
	}
	msgHandler.TaskQueue = make([]chan tiface.IRequest, msgHandler.WorkerPoolSize)
	return msgHandler
}

func (mh *MsgHandler) DoMsgHandler(request tiface.IRequest) {
	msgID := request.GetMsgID()
	handler, ok := mh.Apis[msgID]
	if !ok {
		global.Log.Warnf("%s api msgID=%d is not found", mh.ServerType, msgID)
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router tiface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		global.Log.Warnf("%s api msgID=%d repeated", mh.ServerType, msgID)
		return
	}
	mh.Apis[msgID] = router
	global.Log.Tracef("%s api msgID=%d added", mh.ServerType, msgID)
}

func (mh *MsgHandler) StartOneWorkerPool(workerID int, taskQueue chan tiface.IRequest) {
	global.Log.Tracef("%s worker id=%d started", mh.ServerType, workerID)
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan tiface.IRequest, mh.MaxWorkerTaskLen)
		go mh.StartOneWorkerPool(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request tiface.IRequest) {
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	global.Log.Tracef("%s add connID=%d, request msgID=%d to workerID=%d", mh.ServerType, request.GetConnection().GetConnID(), request.GetMsgID(), workerID)
	mh.TaskQueue[workerID] <- request
}
