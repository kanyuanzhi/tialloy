package tialloy

import (
	"log"
	"testing"
	"tialloy/tiface"
	"tialloy/tinet"
)

type EchoRouter struct {
	tinet.BaseRouter
}

func (ec *EchoRouter) Handle(request tiface.IRequest) {
	log.Printf("[Handle] messageID=%d, %s",request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("server response ECHO"))
	if err != nil {
		log.Println(err.Error())
	}
}

type CustomRouter struct {
	tinet.BaseRouter
}

func (cr *CustomRouter) Handle(request tiface.IRequest) {
	log.Printf("[Handle] messageID=%d, %s",request.GetMsgID(), request.GetData())
	err := request.GetConnection().SendMsg(request.GetMsgID(), []byte("server response CUSTOM"))
	if err != nil {
		log.Println(err.Error())
	}
}

func DoConnStartHook(connection tiface.IConnection)  {
	log.Println("onConnStartHook!!!!!!!!!!!!!!!!!!!!")
	connection.SendBuffMsg(1000, []byte("DoConnStartHook"))
}

func DoConnStopHook(connection tiface.IConnection)  {
	log.Println("onConnStopHook!!!!!!!!!!!!!!!!!!!!")
}

func TestServer(t *testing.T) {
	server := tinet.NewServer("TiAlloy-Test")
	server.AddRouter(1, &CustomRouter{})
	server.AddRouter(2, &EchoRouter{})
	server.SetOnConnStart(DoConnStartHook)
	server.SetOnConnStop(DoConnStopHook)
	server.Serve()
}