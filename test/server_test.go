package test

import (
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/tinet"
	"log"
	"testing"
)

type EchoRouter struct {
	*tinet.BaseRouter
}

func (ec *EchoRouter) Handle(request tiface.IRequest) {
	log.Printf("[Handle] messageID=%d, %s", request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("server response ECHO"))
	if err != nil {
		log.Println(err.Error())
	}
}

type CustomRouter struct {
	*tinet.BaseRouter
}

func (cr *CustomRouter) Handle(request tiface.IRequest) {
	log.Printf("[Handle] messageID=%d, %s", request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("server response CUSTOM"))
	if err != nil {
		log.Println(err.Error())
	}
}

func DoConnStartHook(connection tiface.IConnection) {
	log.Println("onConnStartHook!!!!!!!!!!!!!!!!!!!!")
	connection.SendBuffMsg(1000, []byte("DoConnStartHook"))
}

func DoConnStopHook(connection tiface.IConnection) {
	log.Println("onConnStopHook!!!!!!!!!!!!!!!!!!!!")
}

func TestServer(t *testing.T) {
	server := tinet.NewTcpServer()
	server.AddRouter(1, &CustomRouter{})
	server.AddRouter(2, &EchoRouter{})
	server.SetOnConnStart(DoConnStartHook)
	server.SetOnConnStop(DoConnStopHook)
	server.Serve()
}
