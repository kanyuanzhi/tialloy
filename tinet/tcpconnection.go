package tinet

import (
	"context"
	"errors"
	"fmt"
	"github.com/kanyuanzhi/tialloy/global"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/tilog"
	"io"
	"net"
)

type TcpConnection struct {
	*BaseConnection
}

func NewTcpConnection(server tiface.IServer, conn *net.TCPConn, connID uint32, msgHandler tiface.IMsgHandler) tiface.IConnection {
	baseConnection := NewBaseConnection(server, conn, connID, msgHandler)
	tc := &TcpConnection{
		BaseConnection: baseConnection,
	}
	tc.server.GetConnManager().Add(tc) // 初始化新链接时将新链接加入到链接管理模块中,
	return tc
}

func (tc *TcpConnection) StartReader() {
	tilog.Log.Infof("tcp reader goroutine for %s is running", tc.RemoteAddr())
	defer tilog.Log.Warnf("tcp reader goroutine for %s exited", tc.RemoteAddr())
	defer tc.Stop()

	for {
		select {
		case <-tc.ctx.Done():
			return
		default:
			dp := NewDataPack()

			dataHeadBuf := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(tc.GetTcpConn(), dataHeadBuf); err != nil {
				tilog.Log.Error(err)
				return
			}

			message, err := dp.Unpack(dataHeadBuf)
			if err != nil {
				tilog.Log.Error(err)
				return
			}

			var dataBuf []byte
			if message.GetDataLen() > 0 {
				dataBuf = make([]byte, message.GetDataLen())
				if _, err := io.ReadFull(tc.GetTcpConn(), dataBuf); err != nil {
					tilog.Log.Error(err)
					return
				}
			}

			message.SetData(dataBuf)
			request := NewRequest(tc, message)

			if global.Object.TcpWorkerPoolSize > 0 {
				go tc.MsgHandler.SendMsgToTaskQueue(request)
			} else {
				go tc.MsgHandler.DoMsgHandler(request)
			}
		}
	}
}

func (tc *TcpConnection) StartWriter() {
	tilog.Log.Infof("tcp writer goroutine for %s is running", tc.RemoteAddr())
	defer tilog.Log.Warnf("tcp writer goroutine for %s exited", tc.RemoteAddr())
	for {
		select {
		case data := <-tc.msgChan:
			if _, err := tc.GetTcpConn().Write(data); err != nil {
				tilog.Log.Error(err)
				return
			}
		case data, ok := <-tc.msgBuffChan:
			if ok {
				if _, err := tc.GetTcpConn().Write(data); err != nil {
					tilog.Log.Error(err)
					return
				}
			} else {
				// 通道关闭
				tilog.Log.Error("msgBuffChan has been closed")
				break
			}
		case <-tc.ctx.Done():
			return
		}
	}
}

func (tc *TcpConnection) Start() {
	tc.ctx, tc.cancel = context.WithCancel(context.Background())

	go tc.StartReader()
	go tc.StartWriter()

	tc.server.CallOnConnStart(tc) // 链接启动的回调业务
}

func (tc *TcpConnection) GetTcpConn() *net.TCPConn {
	return tc.Conn.(*net.TCPConn)
}

func (tc *TcpConnection) SendMsg(msgID uint32, data []byte) error {
	tc.RLock()
	defer tc.RUnlock()
	if tc.IsClosed == true {
		return errors.New("tcp connection has been closed")
	}

	dp := NewDataPack()
	binaryMessage, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		return errors.New(fmt.Sprintf("pack tcp messageID=%d err", msgID))
	}

	tc.msgChan <- binaryMessage
	return nil
}

func (tc *TcpConnection) SendBuffMsg(msgID uint32, data []byte) error {
	tc.RLock()
	defer tc.RUnlock()
	if tc.IsClosed == true {
		return errors.New("tcp connection has been closed")
	}

	dp := NewDataPack()
	binaryMessage, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		return errors.New(fmt.Sprintf("pack tcp messageID=%d err", msgID))
	}

	tc.msgBuffChan <- binaryMessage
	return nil
}
