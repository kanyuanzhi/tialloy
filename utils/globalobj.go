package utils

import (
	"encoding/json"
	"gihub.com/kanyuanzhi/tialloy/tiface"
	"io/ioutil"
	"log"
)

type GlobalObj struct {
	TcpServer tiface.IServer `json:"tcp_server,omitempty"`
	Name      string         `json:"name,omitempty"`
	Version   string         `json:"version,omitempty"`
	Host      string         `json:"host,omitempty"`
	TcpPort   int            `json:"tcp_port,omitempty"`

	MaxPacketSize uint32 `json:"max_packet_size,omitempty"`
	MaxConn       int    `json:"max_conn,omitempty"`

	WorkerPoolSize   uint32 `json:"worker_pool_size,omitempty"`    // 工作池大小
	MaxWorkerTaskLen uint32 `json:"max_worker_task_len,omitempty"` // 每个工作池处理的消息队列长度
	MaxMsgChanLen    uint32 `json:"max_msg_chan_len,omitempty"`    // 缓冲数据通道最大缓冲数量
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
	log.Println("[INFO][GlobalObj][Reload] Global object has been reloaded")
}

func init() {
	GlobalObject = &GlobalObj{
		TcpServer: nil,
		Name:      "TiAlloy Server",
		Version:   "V0.1",
		Host:      "127.0.0.1",
		TcpPort:   8888,

		MaxPacketSize: 4096,
		MaxConn:       10000,

		WorkerPoolSize:   20,
		MaxWorkerTaskLen: 10,
		MaxMsgChanLen:    20,
	}

	GlobalObject.Reload()
}
