/*
 * File : main.go
 * CreateDate : 2019-12-07 15:53:28
 * */

package main

import (
	"mzinx/ziface"
	"mzinx/znet"

	"github.com/sirupsen/logrus"
)

type PingRouter struct {
}

func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fun := "PingRouter.PreHandle"
	logrus.Infof("[%s] start.....", fun)
	data := []byte("before ping....")
	_, err := request.GetConnection().GetTCPConnection().Write(data)
	if err != nil {
		logrus.Errorf("[%s] ping failed err:%v", fun, err)
	}
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fun := "PingRouter.Handle"
	logrus.Infof("[%s] start.....", fun)

	data := []byte("ping ping....")
	_, err := request.GetConnection().GetTCPConnection().Write(data)
	if err != nil {
		logrus.Errorf("[%s] ping failed err:%v", fun, err)
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fun := "PingRouter.PostHandle"
	logrus.Infof("[%s] start.....", fun)

	data := []byte("after ping....")
	_, err := request.GetConnection().GetTCPConnection().Write(data)
	if err != nil {
		logrus.Errorf("[%s] ping failed err:%v", fun, err)
	}
}

func main() {
	s := znet.NewServer("[zinx v03]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}

/* vim: set tabstop=4 set shiftwidth=4 */
