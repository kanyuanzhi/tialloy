package tinet

import (
	"github.com/kanyuanzhi/tialloy/global"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/tilog"
)

type BaseServer struct {
	Name       string
	ServerType string // tcp,websocket
	IPVersion  string // tcp4 or other
	IP         string
	Port       int

	msgHandler  tiface.IMsgHandler
	connManager tiface.IConnManager

	OnConnStart func(connection tiface.IConnection)
	OnConnStop  func(connection tiface.IConnection)
}

func NewBaseServer(serverType string) *BaseServer {
	baseServer := &BaseServer{
		Name:       global.Object.Name,
		ServerType: serverType,
		IPVersion:  "tcp4",
		IP:         global.Object.Host,

		msgHandler:  NewMsgHandler(serverType),
		connManager: NewConnManager(),
	}
	switch serverType {
	case "tcp":
		baseServer.Port = global.Object.TcpPort
	case "websocket":
		baseServer.Port = global.Object.WebsocketPort
	}
	return baseServer
}

func (bs *BaseServer) Start() {
	tilog.Log.Panic("implement me")
}

func (bs *BaseServer) Serve() {
	tilog.Log.Panic("implement me")
}

func (bs *BaseServer) Stop() {
	tilog.Log.Warnf("%s server listenner at %s:%d stopped\n", bs.Name, bs.IP, bs.Port)
	bs.connManager.ClearAllConn()
}

func (bs *BaseServer) AddRouter(msgID uint32, router tiface.IRouter) {
	bs.msgHandler.AddRouter(msgID, router)
}

func (bs *BaseServer) GetConnManager() tiface.IConnManager {
	return bs.connManager
}

func (bs *BaseServer) SetOnConnStart(hookFunc func(connection tiface.IConnection)) {
	bs.OnConnStart = hookFunc
}

func (bs *BaseServer) SetOnConnStop(hookFunc func(connection tiface.IConnection)) {
	bs.OnConnStop = hookFunc
}

func (bs *BaseServer) CallOnConnStart(connection tiface.IConnection) {
	if bs.OnConnStart != nil {
		tilog.Log.Tracef("call DoConnStartHook")
		bs.OnConnStart(connection)
	}
}

func (bs *BaseServer) CallOnConnStop(connection tiface.IConnection) {
	if bs.OnConnStop != nil {
		tilog.Log.Tracef("call DoOnConnStopHook")
		bs.OnConnStop(connection)
	}
}

func (bs *BaseServer) GetServerType() string {
	return bs.ServerType
}
