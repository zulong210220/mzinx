package znet

import (
	"io"
	"net"
	"testing"

	"github.com/sirupsen/logrus"
)

// 关闭global
func TestDatapack(t *testing.T) {
	endpoint := "127.0.0.1:7889"
	ln, err := net.Listen("tcp", endpoint)
	if err != nil {
		logrus.Errorf("server ln err:%v", err)
		return
	}

	f := func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				logrus.Errorf("server accept err:%v", err)
				continue
			}

			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						logrus.Errorf("read head err:%v", err)
						return
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						logrus.Errorf("unpack head err:%v", err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							logrus.Errorf("unpack data err:%v", err)
							return
						}

						logrus.Infof("=> Recv Msg Id:%d len:%d data:%s", msg.Id, msg.DataLen, msg.Data)
					}
				}
			}(conn)
		}
	}
	go f()

	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		logrus.Errorf("client dial err:%v", err)
		return
	}

	dp := NewDataPack()
	msg1 := &Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'H', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		logrus.Errorf("client pack msg1 err:%v", err)
		return
	}
	logrus.Infof("sendData :%s", sendData1)

	msg2 := &Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'A', 'e', 'l', 'l', 'o', '-', '-'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		logrus.Errorf("client pack msg2 err:%v", err)
		return
	}

	// 粘包
	sendData1 = append(sendData1, sendData2...)
	logrus.Infof("send data:%s", sendData1)
	conn.Write(sendData1)
	select {}
}

func TestOO(t *testing.T) {
	dp := NewDataPack()
	msg1 := &Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'H', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		logrus.Errorf("client pack msg1 err:%v", err)
		return
	}
	logrus.Infof("sendData :%s", sendData1)

}
