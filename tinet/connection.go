package tinet

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/utils"
	"net"
	"sync"
)

type BaseConnection struct {
	server tiface.IServer // server与connection可以互相索引，使得connection可以通过server.GetConnManager()操作connManager

	Conn interface{}

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

func NewBaseConnection(server tiface.IServer, conn interface{}, connID uint32, msgHandler tiface.IMsgHandler) *BaseConnection {
	baseConnection := &BaseConnection{
		server:     server,
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}
	switch server.GetServerType() {
	case "tcp":
		baseConnection.msgBuffChan = make(chan []byte, utils.GlobalObject.TcpMaxMsgChanLen)
	case "websocket":
		baseConnection.msgBuffChan = make(chan []byte, utils.GlobalObject.WebsocketMaxMsgChanLen)
	}
	return baseConnection
}

func (bc *BaseConnection) Start() {
	panic("implement me")
}

func (bc *BaseConnection) Stop() {
	bc.Lock()
	defer bc.Unlock()

	utils.GlobalLog.Warnf("%s connection connID=%d stopped", bc.server.GetServerType(), bc.ConnID)

	bc.server.CallOnConnStop(bc) //链接关闭的回调业务

	if bc.IsClosed == true {
		return
	}

	switch bc.server.GetServerType() {
	case "tcp":
		bc.Conn.(*net.TCPConn).Close()
	case "websocket":
		bc.Conn.(*websocket.Conn).Close()
	}

	bc.cancel()

	bc.server.GetConnManager().Remove(bc)

	close(bc.msgChan)
	close(bc.msgBuffChan)

	bc.IsClosed = true
}

func (bc *BaseConnection) GetConn() interface{} {
	return bc.Conn
}

func (bc *BaseConnection) GetConnID() uint32 {
	return bc.ConnID
}

func (bc *BaseConnection) RemoteAddr() net.Addr {
	switch bc.server.GetServerType() {
	case "tcp":
		return bc.Conn.(*net.TCPConn).RemoteAddr()
	case "websocket":
		return bc.Conn.(*websocket.Conn).RemoteAddr()
	default:
		return nil
	}
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
