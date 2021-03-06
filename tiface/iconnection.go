package tiface

import (
	"context"
	"github.com/gorilla/websocket"
	"net"
)

type IConnection interface {
	Start()
	Stop()

	GetWebsocketConn() *websocket.Conn
	GetTcpConn() *net.TCPConn

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
