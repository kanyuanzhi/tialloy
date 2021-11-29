package tinet

import (
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

	ExitBuffChan chan bool

	msgChan     chan []byte
	msgBuffChan chan []byte // 带缓冲的数据通道

	property     map[string]interface{}
	propertyLock *sync.RWMutex
}

func NewBaseConnection(server tiface.IServer, conn interface{}, connID uint32, msgHandler tiface.IMsgHandler) *BaseConnection {
	baseConnection := &BaseConnection{
		server:       server,
		Conn:         conn,
		ConnID:       connID,
		IsClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		property:     make(map[string]interface{}),
		propertyLock: new(sync.RWMutex),
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
	utils.GlobalLog.Warnf("%s connection connID=%d stopped", bc.server.GetServerType(), bc.ConnID)
	if bc.IsClosed == true {
		return
	}
	bc.IsClosed = true

	bc.server.CallOnConnStop(bc) //链接关闭的回调业务

	switch bc.server.GetServerType() {
	case "tcp":
		bc.Conn.(*net.TCPConn).Close()
	case "websocket":
		bc.Conn.(*websocket.Conn).Close()
	}

	bc.ExitBuffChan <- true
	close(bc.ExitBuffChan)

	bc.server.GetConnManager().Remove(bc)
	close(bc.msgChan)
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
	bc.propertyLock.RLock()
	defer bc.propertyLock.RUnlock()

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
