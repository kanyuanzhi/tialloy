package tinet

import (
	"errors"
	"fmt"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/utils"
	"io"
	"net"
	"sync"
)

type Connection struct {
	server tiface.IServer // server与connection可以互相索引，使得connection可以通过server.GetConnManager()操作connManager

	Conn       *net.TCPConn
	ConnID     uint32
	IsClosed   bool
	MsgHandler tiface.IMsgHandler

	ExitBuffChan chan bool

	msgChan     chan []byte
	msgBuffChan chan []byte // 带缓冲的数据通道

	property     map[string]interface{}
	propertyLock *sync.RWMutex
}

func NewConnection(server tiface.IServer, conn *net.TCPConn, connID uint32, msgHandler tiface.IMsgHandler) tiface.IConnection {
	c := &Connection{
		server:       server,
		Conn:         conn,
		ConnID:       connID,
		IsClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
		propertyLock: new(sync.RWMutex),
	}
	c.server.GetConnManager().Add(c) // 初始化新链接时将新链接加入到链接管理模块中
	return c
}

func (c *Connection) StartReader() {
	utils.GlobalLog.Infof("reader goroutine for %s is running", c.RemoteAddr())
	defer utils.GlobalLog.Warnf("reader goroutine for %s exited", c.RemoteAddr())
	defer c.Stop()

	for {
		dp := NewDataPack()

		dataHeadBuf := make([]byte, dp.GetHeadLen())

		if _, err := io.ReadFull(c.GetTCPConnection(), dataHeadBuf); err != nil {
			utils.GlobalLog.Error(err)
			c.ExitBuffChan <- true
			return
		}

		message, err := dp.Unpack(dataHeadBuf)
		if err != nil {
			utils.GlobalLog.Error(err)
			c.ExitBuffChan <- true
			return
		}

		var dataBuf []byte
		if message.GetDataLen() > 0 {
			dataBuf = make([]byte, message.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), dataBuf); err != nil {
				utils.GlobalLog.Error(err)
				continue
			}
		}

		message.SetData(dataBuf)
		request := NewRequest(c, message)

		if utils.GlobalObject.WorkerPoolSize > 0 {
			go c.MsgHandler.SendMsgToTaskQueue(request)
		} else {
			go c.MsgHandler.DoMsgHandler(request)
		}
	}
}

func (c *Connection) StartWriter() {
	utils.GlobalLog.Infof("writer goroutine for %s is running", c.RemoteAddr())
	defer utils.GlobalLog.Warnf("writer goroutine for %s exited", c.RemoteAddr())
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				utils.GlobalLog.Error(err)
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					utils.GlobalLog.Error(err)
					return
				}
			} else {
				// 通道关闭
				utils.GlobalLog.Error("msgBuffChan has been closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()

	c.server.CallOnConnStart(c) // 链接启动的回调业务

	for {
		select {
		case <-c.ExitBuffChan: // 得到退出消息，不再阻塞
			return
		}
	}
}

func (c *Connection) Stop() {
	utils.GlobalLog.Warnf("connection connID=%d stopped", c.ConnID)
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	c.server.CallOnConnStop(c) //链接关闭的回调业务

	c.Conn.Close()
	c.ExitBuffChan <- true
	close(c.ExitBuffChan)

	c.server.GetConnManager().Remove(c)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(messageID uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("connection has been closed")
	}

	dp := NewDataPack()
	binaryMessage, err := dp.Pack(NewMessage(messageID, data))
	if err != nil {
		return errors.New(fmt.Sprintf("pack messageID=%d err", messageID))
	}

	c.msgChan <- binaryMessage

	return nil
}

func (c *Connection) SendBuffMsg(messageID uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("connection has been closed")
	}

	dp := NewDataPack()
	binaryMessage, err := dp.Pack(NewMessage(messageID, data))
	if err != nil {
		return errors.New(fmt.Sprintf("pack messageID=%d err", messageID))
	}

	c.msgBuffChan <- binaryMessage

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
