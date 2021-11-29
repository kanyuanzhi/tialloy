package tiface

import (
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
}
