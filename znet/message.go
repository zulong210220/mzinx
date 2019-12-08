package znet

import "fmt"

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(dl uint32) {
	m.DataLen = dl
}

func (m *Message) String() string {
	if m == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{id:%d,len:%d,data:'%s'}", m.Id, m.DataLen, m.Data)
}
