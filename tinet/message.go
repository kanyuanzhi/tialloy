package tinet

import (
	"github.com/kanyuanzhi/tialloy/tiface"
)

type Message struct {
	DataLen uint32
	MsgID   uint32
	Data    []byte
}

func NewMessage(msgID uint32, data []byte) tiface.IMessage {
	return &Message{
		DataLen: uint32(len(data)),
		MsgID:   msgID,
		Data:    data,
	}
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}

func (m *Message) SetMsgID(msgID uint32) {
	m.MsgID = msgID
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
