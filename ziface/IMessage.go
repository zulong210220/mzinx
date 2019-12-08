package ziface

type IMessage interface {
	GetMsgId() uint32
	GetDataLen() uint32
	GetData() []byte

	SetMsgId(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
