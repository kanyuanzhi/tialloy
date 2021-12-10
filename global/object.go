package global

import (
	"encoding/json"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

var Object *Obj

type Obj struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Host    string `json:"host,omitempty"`

	// tcp settings
	TcpServer           tiface.IServer `json:"tcp_server,omitempty"`
	TcpPort             int            `json:"tcp_port,omitempty"`
	TcpMaxPacketSize    uint32         `json:"tcp_max_packet_size,omitempty"`
	TcpMaxConn          int            `json:"tcp_max_conn,omitempty"`
	TcpWorkerPoolSize   uint32         `json:"tcp_worker_pool_size,omitempty"`    // 工作池大小
	TcpMaxWorkerTaskLen uint32         `json:"tcp_max_worker_task_len,omitempty"` // 每个工作池处理的消息队列长度
	TcpMaxMsgChanLen    uint32         `json:"tcp_max_msg_chan_len,omitempty"`    // 缓冲数据通道最大缓冲数量

	// websocket settings
	WebsocketServer           tiface.IServer `json:"websocket_server,omitempty"`
	WebsocketPort             int            `json:"websocket_port,omitempty"`
	WebsocketScheme           string         `json:"websocket_scheme,omitempty"` // websocket协议，ws、wss
	WebsocketPath             string         `json:"websocket_path,omitempty"`   // websocket请求路径
	WebsocketMaxConn          int            `json:"websocket_max_conn,omitempty"`
	WebsocketWorkerPoolSize   uint32         `json:"websocket_worker_pool_size,omitempty"`
	WebsocketMaxWorkerTaskLen uint32         `json:"websocket_max_worker_task_len,omitempty"`
	WebsocketMaxMsgChanLen    uint32         `json:"websocket_max_msg_chan_len,omitempty"`

	LogMode bool `json:"log_mode,omitempty"`
}

func (g *Obj) Reload() {
	data, err := ioutil.ReadFile("conf/tialloy.json")
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(data, &Object)
	if err != nil {
		panic(err.Error())
	}

	Log = logrus.New()
	Log.SetReportCaller(Object.LogMode)
	if Object.LogMode == true {
		Log.SetLevel(logrus.TraceLevel)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}
	Log.SetFormatter(&customFormatter{})
}

func init() {
	Object = &Obj{
		Name:    "TiAlloy Server",
		Version: "V0.1",
		Host:    "127.0.0.1",

		TcpServer:           nil,
		TcpPort:             8888,
		TcpMaxPacketSize:    4096,
		TcpMaxConn:          1000,
		TcpWorkerPoolSize:   20,
		TcpMaxWorkerTaskLen: 10,
		TcpMaxMsgChanLen:    20,

		WebsocketServer:           nil,
		WebsocketPort:             10000,
		WebsocketPath:             "touch",
		WebsocketMaxConn:          10000,
		WebsocketWorkerPoolSize:   20,
		WebsocketMaxWorkerTaskLen: 10,
		WebsocketMaxMsgChanLen:    20,

		LogMode: true, // true：详细，打印log在代码中输出位置；false：简要，不打印文件输出位置，不打印debug和trace（性能高，生产环境使用）
	}

	Object.Reload()
}
