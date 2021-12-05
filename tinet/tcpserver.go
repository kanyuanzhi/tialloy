package tinet

import (
	"fmt"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/utils"
	"math/rand"
	"net"
	"time"
)

type TcpServer struct {
	*BaseServer
}

func NewTcpServer() tiface.IServer {
	baseServer := NewBaseServer("tcp")
	return &TcpServer{
		BaseServer: baseServer,
	}
}

func (ts *TcpServer) Start() {
	utils.GlobalLog.Infof("%s tcp server listenner on %s:%d is starting...", ts.Name, ts.IP, ts.Port)
	go func() {
		ts.msgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(ts.IPVersion, fmt.Sprintf("%s:%d", ts.IP, ts.Port))
		if err != nil {
			utils.GlobalLog.Error(err)
			return
		}
		listener, err := net.ListenTCP(ts.IPVersion, addr)
		if err != nil {
			utils.GlobalLog.Error(err)
			return
		}
		utils.GlobalLog.Infof("tcp server is listening on %s:%d", ts.IP, ts.Port)

		for {
			conn, err := listener.AcceptTCP() // 阻塞等待客户端建立连接请求
			if err != nil {
				utils.GlobalLog.Error(err)
				return
			}

			if ts.connManager.Len() >= utils.GlobalObject.TcpMaxConn {
				// TODO:此处应通知客户端服务器拒绝服务?
				conn.Close() // 超过服务器设置的最大TCP连接数，拒绝服务
				continue
			}

			rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
			connID := rand.Uint32()

			dealConn := NewTcpConnection(ts, conn, connID, ts.msgHandler)

			go dealConn.Start()
		}
	}()
}

func (ts *TcpServer) Serve() {
	ts.Start()
	//TODO 服务器启动后的一些操作
	select {}
}
