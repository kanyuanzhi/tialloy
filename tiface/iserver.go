package tiface

type IServer interface {
	Start()
	Stop()
	Serve()

	AddRouter(msgID uint32, router IRouter)
	GetConnManager() IConnManager

	SetOnConnStart(func(connection IConnection))
	SetOnConnStop(func(connection IConnection))
	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)
}
