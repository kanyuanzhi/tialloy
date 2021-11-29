package tinet

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/utils"
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
	utils.GlobalLog.Infof("websocket reader goroutine for %s is running", wc.RemoteAddr())
	defer utils.GlobalLog.Warnf("websocket reader goroutine for %s exited", wc.RemoteAddr())
	defer wc.Stop()

	for {
		msgType, data, err := wc.Conn.(*websocket.Conn).ReadMessage()
		if err != nil {
			utils.GlobalLog.Error(err)
			wc.ExitBuffChan <- true
			return
		}
		wc.MessageType = msgType

		var msgJon map[string]interface{}
		if err := json.Unmarshal(data, &msgJon); err != nil {
			utils.GlobalLog.Error(err)
			wc.ExitBuffChan <- true
			return
		}

		if msgID, ok := msgJon["msgID"]; ok {
			message := NewMessage(uint32(msgID.(float64)), data)
			request := NewRequest(wc, message)

			if utils.GlobalObject.WorkerPoolSize > 0 {
				go wc.MsgHandler.SendMsgToTaskQueue(request)
			} else {
				go wc.MsgHandler.DoMsgHandler(request)
			}
		} else {
			utils.GlobalLog.Warn("no msgID")
		}
	}
}

func (wc *WebsocketConnection) StartWriter() {
	utils.GlobalLog.Infof("websocket writer goroutine for %s is running", wc.RemoteAddr())
	defer utils.GlobalLog.Warnf("websocket writer goroutine for %s exited", wc.RemoteAddr())
	for {
		select {
		case data := <-wc.msgChan:
			if err := wc.Conn.(*websocket.Conn).WriteMessage(wc.MessageType, data); err != nil {
				utils.GlobalLog.Error(err)
				return
			}
		case data, ok := <-wc.msgBuffChan:
			if ok {
				if err := wc.Conn.(*websocket.Conn).WriteMessage(wc.MessageType, data); err != nil {
					utils.GlobalLog.Error(err)
					return
				}
			} else {
				// 通道关闭
				utils.GlobalLog.Error("msgBuffChan has been closed")
				break
			}
		case <-wc.ExitBuffChan:
			return
		}
	}
}

func (wc *WebsocketConnection) Start() {
	go wc.StartReader()
	go wc.StartWriter()

	wc.server.CallOnConnStart(wc)

	for {
		select {
		case <-wc.ExitBuffChan: // 得到退出消息，不再阻塞
			return
		}
	}
}

func (wc *WebsocketConnection) GetWebsocketConn() *websocket.Conn {
	return wc.Conn.(*websocket.Conn)
}

func (wc *WebsocketConnection) SendMsg(msgID uint32, data []byte) error {
	if wc.IsClosed == true {
		return errors.New("websocket connection has been closed")
	}
	wc.msgChan <- data
	return nil
}

func (wc *WebsocketConnection) SendBuffMsg(msgID uint32, data []byte) error {
	if wc.IsClosed == true {
		return errors.New("websocket connection has been closed")
	}
	wc.msgBuffChan <- data
	return nil
}
