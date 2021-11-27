package tinet

import (
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/utils"
	"log"
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
		log.Printf("api msgID=%d is not found", msgID)
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router tiface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		log.Printf("api msgID=%d repeated", msgID)
		return
	}
	mh.Apis[msgID] = router
	log.Printf("api msgID=%d is added", msgID)
}

func (mh *MsgHandler) StartOneWorkerPool(workerID int, taskQueue chan tiface.IRequest) {
	log.Printf("worker id=%d is started", workerID)
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
	log.Printf("add connID=%d, request msgID=%d to workerID=%d", request.GetConnection().GetConnID(), request.GetMsgID(), workerID)
	mh.TaskQueue[workerID] <- request
}
