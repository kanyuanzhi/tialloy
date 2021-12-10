package tinet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kanyuanzhi/tialloy/global"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/tilog"
	"math/rand"
	"net/http"
	"time"
)

type WebsocketServer struct {
	*BaseServer
	Scheme string
	Path   string
}

func NewWebsocketServer() tiface.IServer {
	baseServer := NewBaseServer("websocket")
	websocketServer := &WebsocketServer{
		BaseServer: baseServer,
		Scheme:     "ws",
		Path:       global.Object.WebsocketPath,
	}
	return websocketServer
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  int(global.Object.TcpMaxPacketSize), //读取最大值
	WriteBufferSize: int(global.Object.TcpMaxPacketSize), //写最大值
	//解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *WebsocketServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		tilog.Log.Error(err)
		return
	}

	if ws.connManager.Len() >= global.Object.TcpMaxConn {
		// TODO:此处应通知客户端服务器拒绝服务?
		tilog.Log.Warnf("connection num reaches max %d", global.Object.TcpMaxConn)
		conn.Close() // 超过服务器设置的最大TCP连接数，拒绝服务
		return
	}

	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	connID := rand.Uint32()

	dealConn := NewWebsocketConnection(ws, conn, connID, ws.msgHandler)

	go dealConn.Start()

}

func (ws *WebsocketServer) Start() {
	tilog.Log.Infof("%s websocket server listenner on %s:%d is starting...", ws.Name, ws.IP, ws.Port)
	ws.msgHandler.StartWorkerPool()
	http.HandleFunc("/"+ws.Path, ws.wsHandler)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ws.IP, ws.Port), nil)
	if err != nil {
		tilog.Log.Error(err)
	}
}

func (ws *WebsocketServer) Serve() {
	ws.Start()
	//TODO 服务器启动后的一些操作
	select {}
}
