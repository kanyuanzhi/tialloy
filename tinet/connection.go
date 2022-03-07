package tinet

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/kanyuanzhi/tialloy/global"
	"github.com/kanyuanzhi/tialloy/tiface"
	"net"
	"sync"
)

type BaseConnection struct {
	server tiface.IServer // server与connection可以互相索引，使得connection可以通过server.GetConnManager()操作connManager

	//Conn interface{}

	ConnID     uint32
	IsClosed   bool
	MsgHandler tiface.IMsgHandler

	ctx    context.Context
	cancel context.CancelFunc

	msgChan     chan []byte
	msgBuffChan chan []byte // 带缓冲的数据通道

	sync.RWMutex

	property     map[string]interface{}
	propertyLock sync.Mutex
}

func NewBaseConnection(server tiface.IServer, connID uint32, msgHandler tiface.IMsgHandler) *BaseConnection {
	baseConnection := &BaseConnection{
		server:     server,
		//Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}
	switch server.GetServerType() {
	case "tcp":
		baseConnection.msgBuffChan = make(chan []byte, global.Object.TcpMaxMsgChanLen)
	case "websocket":
		baseConnection.msgBuffChan = make(chan []byte, global.Object.WebsocketMaxMsgChanLen)
	}
	return baseConnection
}

func (bc *BaseConnection) Start() {
	panic("implement me")
}

func (bc *BaseConnection) Stop() {
	panic("implement me")
}

func (bc *BaseConnection) GetTcpConn() *net.TCPConn {
	panic("implement me")
}

func (bc *BaseConnection) GetWebsocketConn() *websocket.Conn {
	panic("implement me")
}

func (bc *BaseConnection) GetConnID() uint32 {
	return bc.ConnID
}

func (bc *BaseConnection) RemoteAddr() net.Addr {
	panic("implement me")
}

func (bc *BaseConnection) SendMsg(msgID uint32, data []byte) error {
	panic("implement me")
}

func (bc *BaseConnection) SendBuffMsg(msgID uint32, data []byte) error {
	panic("implement me")
}

func (bc *BaseConnection) SetProperty(key string, value interface{}) {
	bc.propertyLock.Lock()
	defer bc.propertyLock.Unlock()

	bc.property[key] = value
}

func (bc *BaseConnection) GetProperty(key string) (interface{}, error) {
	bc.propertyLock.Lock()
	defer bc.propertyLock.Unlock()

	if value, ok := bc.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (bc *BaseConnection) RemoveProperty(key string) {
	bc.propertyLock.Lock()
	defer bc.propertyLock.Unlock()

	delete(bc.property, key)
}

func (bc *BaseConnection) GetServer() tiface.IServer {
	return bc.server
}

func (bc *BaseConnection) Context() context.Context {
	return bc.ctx
}
