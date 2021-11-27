package tinet

import "gihub.com/kanyuanzhi/tialloy/tiface"

type Request struct {
	conn    tiface.IConnection
	message tiface.IMessage
}

func (r *Request) GetConnection() tiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.message.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.message.GetMsgID()
}

func NewRequest(conn tiface.IConnection, message tiface.IMessage) tiface.IRequest {
	return &Request{
		conn:    conn,
		message: message,
	}
}
