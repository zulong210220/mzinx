/*
 * File : client.go
 * CreateDate : 2019-12-07 16:17:25
 * */
package main

import (
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

	bufSize := 1024
	for {
		buf := []byte("Hello zinx v01---")
		_, err := conn.Write(buf)
		if err != nil {
			logrus.Errorf("conn read buf:%s failed err:%v", buf, err)
			return
		}

		readBuf := make([]byte, bufSize)
		cnt, err := conn.Read(readBuf)
		if err != nil {
			logrus.Errorf("read buf failed err:%v", err)
			return
		}

		logrus.Infof("server callback %s cnt:%d", buf, cnt)
		time.Sleep(10 * time.Millisecond)
	}
}

/* vim: set tabstop=4 set shiftwidth=4 */
