package tiface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgID uint32, router IRouter)

	StartWorkerPool() // 启动线程池
	SendMsgToTaskQueue(request IRequest) // 将请求交给任务队列
}




