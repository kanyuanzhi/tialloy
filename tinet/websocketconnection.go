package tinet

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/kanyuanzhi/tialloy/global"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/tilog"
)

type WebsocketConnection struct {
	*BaseConnection
	MessageType int
}

func NewWebsocketConnection(server tiface.IServer, conn *websocket.Conn, connID uint32, msgHandler tiface.IMsgHandler) tiface.IConnection {
	baseConnection := NewBaseConnection(server, conn, connID, msgHandler)
	wc := &WebsocketConnection{
		BaseConnection: baseConnection,
		MessageType:    websocket.TextMessage, // 默认文本协议
	}
	wc.server.GetConnManager().Add(wc) // 初始化新链接时将新链接加入到链接管理模块中,
	return wc
}

func (wc *WebsocketConnection) StartReader() {
	tilog.Log.Infof("websocket reader goroutine for %s is running", wc.RemoteAddr())
	defer tilog.Log.Warnf("websocket reader goroutine for %s exited", wc.RemoteAddr())
	defer wc.Stop()

	for {
		select {
		case <-wc.ctx.Done():
			return
		default:
			msgType, data, err := wc.GetWebsocketConn().ReadMessage()
			if err != nil {
				tilog.Log.Error(err)
				return
			}
			wc.MessageType = msgType

			var msgJon map[string]interface{}
			if err := json.Unmarshal(data, &msgJon); err != nil {
				tilog.Log.Error(err)
				return
			}
			if msgID, ok := msgJon["msg_id"]; ok {
				message := NewMessage(uint32(msgID.(float64)), data)
				request := NewRequest(wc, message)
				if global.Object.WebsocketWorkerPoolSize > 0 {
					go wc.MsgHandler.SendMsgToTaskQueue(request)
				} else {
					go wc.MsgHandler.DoMsgHandler(request)
				}
			} else {
				tilog.Log.Warn("no msg_id")
			}
		}
	}
}

func (wc *WebsocketConnection) StartWriter() {
	tilog.Log.Infof("websocket writer goroutine for %s is running", wc.RemoteAddr())
	defer tilog.Log.Warnf("websocket writer goroutine for %s exited", wc.RemoteAddr())
	for {
		select {
		case data := <-wc.msgChan:
			if err := wc.GetWebsocketConn().WriteMessage(wc.MessageType, data); err != nil {
				tilog.Log.Error(err)
				return
			}
		case data, ok := <-wc.msgBuffChan:
			if ok {
				if err := wc.GetWebsocketConn().WriteMessage(wc.MessageType, data); err != nil {
					tilog.Log.Error(err)
					return
				}
			} else {
				// 通道关闭
				tilog.Log.Error("msgBuffChan has been closed")
				break
			}
		case <-wc.ctx.Done():
			return
		}
	}
}

func (wc *WebsocketConnection) Start() {
	wc.ctx, wc.cancel = context.WithCancel(context.Background())

	go wc.StartReader()
	go wc.StartWriter()

	wc.server.CallOnConnStart(wc)
}

func (wc *WebsocketConnection) GetWebsocketConn() *websocket.Conn {
	return wc.Conn.(*websocket.Conn)
}

func (wc *WebsocketConnection) SendMsg(msgID uint32, data []byte) error {
	wc.RLock()
	defer wc.RUnlock()
	if wc.IsClosed == true {
		return errors.New("websocket connection has been closed")
	}
	wc.msgChan <- data
	return nil
}

func (wc *WebsocketConnection) SendBuffMsg(msgID uint32, data []byte) error {
	wc.RLock()
	defer wc.RUnlock()
	if wc.IsClosed == true {
		return errors.New("websocket connection has been closed")
	}
	wc.msgBuffChan <- data
	return nil
}
