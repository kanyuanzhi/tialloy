package tinet

import (
	"fmt"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/utils"
	log "github.com/sirupsen/logrus"
	//"log"
	"math/rand"
	"net"
	"time"
)

type Server struct {
	Name      string
	IPVersion string //tcp4 or other
	IP        string
	Port      int

	msgHandler  tiface.IMsgHandler
	connManager tiface.IConnManager

	OnConnStart func(connection tiface.IConnection)
	OnConnStop  func(connection tiface.IConnection)
}

func (s *Server) Start() {
	//utils.GlobalLog.Info("123")
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
	log.Info(123214)
	//customFormatter := new(logrus.TextFormatter)
	//customFormatter.TimestampFormat =  "2006-01-02 15:04:05"
	//logrus.SetFormatter(customFormatter)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	//customFormatter.FullTimestamp = true
	log.Info(123214)

	log.Printf("[INFO][Server][START] Server listenner at %s:%d, is starting...\n", s.IP, s.Port)

	go func() {

		s.msgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			log.Println("[ERROR][Server][START][ResolveTCPAddr]", err.Error())
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Println("[ERROR][Server][START][ListenTCP]", err.Error())
			return
		}
		log.Printf("[INFO][Server][START] The tiAlloy server %s is listening on %s:%d", s.Name, s.IP, s.Port)

		for {
			conn, err := listener.AcceptTCP() // 阻塞等待客户端建立连接请求
			if err != nil {
				log.Println("[ERROR][Server][START][AcceptTCP]", err.Error())
				return
			}

			if s.connManager.Len() >= utils.GlobalObject.MaxConn {
				// TODO:此处应通知客户端服务器拒绝服务?
				conn.Close() // 超过服务器设置的最大TCP连接数，拒绝服务
				continue
			}

			connID := rand.Uint32()

			dealConn := NewConnection(s, conn, connID, s.msgHandler)

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	//TODO 清理连接
	fmt.Printf("[INFO][Server][STOP] Server listenner at IP: %s, Port %d, is stopped\n", s.IP, s.Port)
	s.connManager.ClearAllConn()
}

func (s *Server) Serve() {
	s.Start()

	//TODO 服务器启动后的一些操作

	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(msgID uint32, router tiface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnManager() tiface.IConnManager {
	return s.connManager
}

func (s *Server) SetOnConnStart(hookFunc func(connection tiface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(connection tiface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(connection tiface.IConnection) {
	if s.OnConnStart != nil {
		log.Println("CallOnConnStart->")
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnStop(connection tiface.IConnection) {
	if s.OnConnStop != nil {
		log.Println("CallOnConnStop->")
		s.OnConnStop(connection)
	}
}

func NewServer(name string) tiface.IServer {
	//utils.GlobalObject.Reload()
	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,

		msgHandler:  NewMsgHandler(),
		connManager: NewConnManager(),
	}
}
