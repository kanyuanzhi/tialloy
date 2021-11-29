package tiface

type IMessage interface {
	GetDataLen() uint32
	GetMsgID() uint32
	GetData() []byte

	SetDataLen(dataLen uint32)
	SetMsgID(msgID uint32)
	SetData(data []byte)
}
