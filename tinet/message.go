package tinet

type Message struct {
	MsgID uint32
	Data  []byte
	DataLen   uint32
}

func NewMessage(msgID uint32, data []byte) *Message {
	return &Message{
		MsgID:   msgID,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) SetMsgID(msgID uint32) {
	m.MsgID = msgID
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}
