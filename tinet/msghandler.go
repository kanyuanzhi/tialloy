package tinet

import (
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/utils"
)

type MsgHandler struct {
	Apis           map[uint32]tiface.IRouter // Apis[msgID] = handler
	WorkerPoolSize uint32
	TaskQueue      []chan tiface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]tiface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan tiface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandler) DoMsgHandler(request tiface.IRequest) {
	msgID := request.GetMsgID()
	handler, ok := mh.Apis[msgID]
	if !ok {
		utils.GlobalLog.Warnf("api msgID=%d is not found", msgID)
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router tiface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		utils.GlobalLog.Warnf("api msgID=%d repeated", msgID)
		return
	}
	mh.Apis[msgID] = router
	utils.GlobalLog.Tracef("api msgID=%d added", msgID)
}

func (mh *MsgHandler) StartOneWorkerPool(workerID int, taskQueue chan tiface.IRequest) {
	utils.GlobalLog.Tracef("worker id=%d started", workerID)
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan tiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorkerPool(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request tiface.IRequest) {
	workerID := request.GetMsgID() % mh.WorkerPoolSize
	utils.GlobalLog.Tracef("add connID=%d, request msgID=%d to workerID=%d", request.GetConnection().GetConnID(), request.GetMsgID(), workerID)
	mh.TaskQueue[workerID] <- request
}
