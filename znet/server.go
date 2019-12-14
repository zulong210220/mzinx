/*
 * File : server.go
 * CreateDate : 2019-12-07 15:21:46
 * */

package znet

import (
	"fmt"
	"mzinx/ziface"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mzinx/utils"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int32
	connCount uint32

	handler     ziface.IMsgHandler
	connManager ziface.IConnManager

	onConnStart func(conn ziface.IConnection)
	onConnStop  func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fun := "Server.Start"
	endpoint := fmt.Sprintf("%s:%d", s.IP, s.Port)
	logrus.Infof("[%s] Server:%s Listener IP:%s Port:%d version:%s starting...", fun, s.Name, s.IP, s.Port, utils.GlobalObject.Version)

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

		if s.connManager.Len() > utils.GlobalObject.MaxConn {
			conn.Close()
			continue
		}

		dealConn := NewConnection(s, conn, s.connCount, s.handler)
		s.connCount++
		go dealConn.Start()
	}
}

func Callback2Client(conn *net.TCPConn, data []byte, cnt int) error {
	fun := "Callback2Client"
	if _, err := conn.Write(data[:cnt]); err != nil {
		logrus.Errorf("[%s] Write data:%s failed err:%v", fun, data, err)
		return err
	}
	return nil
}

func (s *Server) Stop() {
	s.connManager.ClearConn()
	s.handler.Stop()
}

func (s *Server) Serve() {
	go s.Start()
	defer s.Stop()

	// 其它初始化

	// 阻塞
	signalProc()
}

func signalProc() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGALRM, syscall.SIGTERM, syscall.SIGUSR1)

	sig := <-c

	logrus.Warnf("Signal received: %v", sig)

	time.Sleep(100 * time.Millisecond)
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	fun := "Server.AddRouter"
	logrus.Infof("[%s] router:%v", fun, router)
	s.handler.AddRouter(msgId, router)
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		handler:     NewMsgHandler(),
		connManager: NewConnManager(),
	}

	return s
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.connManager
}

func (s *Server) SetOnConnStart(hook func(ziface.IConnection)) {
	s.onConnStart = hook
}

func (s *Server) SetOnConnStop(hook func(ziface.IConnection)) {
	s.onConnStop = hook
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.onConnStart != nil {
		s.onConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.onConnStop != nil {
		s.onConnStop(conn)
	}
}

/* vim: set tabstop=4 set shiftwidth=4 */
