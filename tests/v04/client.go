/*
 * File : client.go
 * CreateDate : 2019-12-07 16:17:25
 * */
package main

import (
	"io"
	"mzinx/znet"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infof("client start...")

	time.Sleep(10 * time.Millisecond)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		logrus.Errorf("client start failed err:%v", err)
		return
	}

	for {
		buf := []byte("Hello zinx v04---")
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPack(0, buf))
		_, err := conn.Write(msg)
		if err != nil {
			logrus.Errorf("write buf:%s error:%v", buf, err)
			return
		}

		readBuf := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, readBuf)
		if err != nil {
			logrus.Errorf("conn read head failed err:%v", err)
			return
		}

		logrus.Infof("pppppp [%s]", readBuf)
		msgHead, err := dp.Unpack(readBuf)
		if err != nil {
			logrus.Errorf("unpack failed err:%v", err)
			return
		}

		if msgHead.GetDataLen() <= 0 {
			continue
		}

		msgRecv := msgHead.(*znet.Message)
		msgRecv.Data = make([]byte, msgRecv.GetDataLen())

		_, err = io.ReadFull(conn, msgRecv.Data)
		if err != nil {
			logrus.Errorf("server unpack data err:%v", err)
			return
		}

		logrus.Infof("Recv Msg Id:%d len:%d data:%s", msgRecv.Id, msgRecv.DataLen, msgRecv.Data)
		time.Sleep(3 * time.Second)
	}
}

/* vim: set tabstop=4 set shiftwidth=4 */
