package utils

import (
	"encoding/json"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type GlobalObj struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Host    string `json:"host,omitempty"`

	TcpServer        tiface.IServer `json:"tcp_server,omitempty"`
	TcpPort          int            `json:"tcp_port,omitempty"`
	MaxPacketSize    uint32         `json:"max_packet_size,omitempty"`
	MaxConn          int            `json:"max_conn,omitempty"`
	WorkerPoolSize   uint32         `json:"worker_pool_size,omitempty"`    // 工作池大小
	MaxWorkerTaskLen uint32         `json:"max_worker_task_len,omitempty"` // 每个工作池处理的消息队列长度
	MaxMsgChanLen    uint32         `json:"max_msg_chan_len,omitempty"`    // 缓冲数据通道最大缓冲数量

	WebsocketServer         tiface.IServer `json:"websocket_server,omitempty"`
	WebsocketPort           int            `json:"websocket_port,omitempty"`
	WebsocketPath           string         `json:"websocket_path,omitempty"`
	WebsocketMaxConn        int            `json:"websocket_max_conn,omitempty"`
	WebsocketWorkerPoolSize uint32         `json:"websocket_worker_pool_size,omitempty"` // 工作池大小

	LogMode bool `json:"log_mode,omitempty"`
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/tialloy.json")
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err.Error())
	}

	GlobalLog = logrus.New()
	GlobalLog.SetReportCaller(GlobalObject.LogMode)
	if GlobalObject.LogMode == true {
		GlobalLog.SetLevel(logrus.TraceLevel)
	} else {
		GlobalLog.SetLevel(logrus.InfoLevel)
	}
	GlobalLog.SetFormatter(&customFormatter{})
}

func init() {
	GlobalObject = &GlobalObj{
		Name:    "TiAlloy Server",
		Version: "V0.1",
		Host:    "127.0.0.1",

		TcpServer:        nil,
		TcpPort:          8888,
		MaxPacketSize:    4096,
		MaxConn:          1000,
		WorkerPoolSize:   20,
		MaxWorkerTaskLen: 10,
		MaxMsgChanLen:    20,

		WebsocketServer:         nil,
		WebsocketPort:           10000,
		WebsocketPath:           "touch",
		WebsocketMaxConn:        10000,
		WebsocketWorkerPoolSize: 20,

		LogMode: true, // true：详细，打印log输出位置；false：简要，不打印文件输出位置，不打印debug和trace（性能高，生产环境使用）
	}

	GlobalObject.Reload()
}
