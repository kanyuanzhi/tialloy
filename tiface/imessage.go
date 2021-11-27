package tiface

type IMessage interface {
	GetMsgID() uint32
	GetData() []byte
	GetDataLen() uint32

	SetMsgID(messageID uint32)
	SetData(data []byte)
	SetDataLen(dataLen uint32)
}
