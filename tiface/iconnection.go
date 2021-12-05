package tiface

import (
	"context"
	"net"
)

type IConnection interface {
	Start()
	Stop()

	GetConn() interface{}

	GetConnID() uint32
	RemoteAddr() net.Addr

	SendMsg(msgID uint32, data []byte) error
	SendBuffMsg(msgID uint32, data []byte) error

	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)

	GetServer() IServer

	Context() context.Context // 用于用户自定义的go程获取连接退出的状态
}
