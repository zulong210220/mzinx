package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"mzinx/consts"
	"mzinx/utils"
	"mzinx/ziface"

	"github.com/sirupsen/logrus"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return consts.DefaultHeadLen
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	fun := "DataPack.Pack"
	dataBuf := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		logrus.Errorf("[%s] binary Len Write data:%s err:%v", fun, msg, err)
		return nil, err
	}

	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		logrus.Errorf("[%s] binary Id Write data:%s err:%v", fun, msg, err)
		return nil, err
	}

	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		logrus.Errorf("[%s] binary Data Write data:%s err:%v", fun, msg, err)
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	fun := "DataPack.Unpack"
	dataBuf := bytes.NewReader(data)

	msg := &Message{}

	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		logrus.Errorf("[%s] binary Read Len err:%s", fun, err)
		return nil, err
	}

	if utils.GlobalObject.MaxPackageSize > 0 &&
		msg.DataLen > utils.GlobalObject.MaxPackageSize {
		logrus.Errorf("[%s] binary Read Len %d > %d", fun, msg.DataLen, utils.GlobalObject.MaxPackageSize)
		return nil, errors.New("over size")
	}

	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		logrus.Errorf("[%s] binary Read Id err:%s", fun, err)
		return nil, err
	}

	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Data); err != nil {
		logrus.Errorf("[%s] binary Read Data err:%s", fun, err)
		return nil, err
	}

	return msg, nil
}
