/*
 * File : server.go
 * CreateDate : 2019-12-07 15:21:46
 * */

package znet

import (
	"fmt"
	"io"
	"mzinx/ziface"
	"net"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

const (
	LineBufferSize = 1024
)

func (s *Server) Start() {
	fun := "Server.Start"
	endpoint := fmt.Sprintf("%s:%d", s.IP, s.Port)
	logrus.Infof("[%s] Server Listener IP:%s Port:%d starting...", fun, s.IP, s.Port)

	addr, err := net.ResolveTCPAddr(s.IPVersion, endpoint)
	if err != nil {
		logrus.Errorf("resolve tcp addr err:%v", err)
		return
	}

	ln, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		logrus.Errorf("[%s] listen %s endpoint:%s err:%v", fun, s.IPVersion, endpoint, err)
		return
	}

	logrus.Infof("[%s] start zinx server success %s %s listening", fun, s.Name, endpoint)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			logrus.Errorf("[%s] AcceptTCP %s err:%s", fun, endpoint, err)
			continue
		}

		go func() {
			for {
				buf := make([]byte, LineBufferSize)
				cnt, err := conn.Read(buf)
				if err != nil {
					logrus.Errorf("[%s] recv buffer %s err:%v", fun, endpoint, err)
					if err == io.EOF {
						break
					}
					continue
				}

				// 回显
				if _, err := conn.Write(buf[:cnt]); err != nil {
					logrus.Errorf("[%s] write back buffer %s %s err:%v", fun, endpoint, buf, err)
					continue
				}
			}
		}()
	}
}

func (s *Server) Stop() {
}

func (s *Server) Serve() {
	go s.Start()

	// 其它初始化

	// 阻塞
	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}

/* vim: set tabstop=4 set shiftwidth=4 */
