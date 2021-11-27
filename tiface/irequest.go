package tiface

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte //获取请求消息的数据
	GetMsgID() uint32
}
