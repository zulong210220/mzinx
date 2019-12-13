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
	logrus.Infof("[%s] start.....msgId:%d msg:%s", fun, request.GetMsgId(), request.GetData())
	/*
		data := []byte("before ping....")
		_, err := request.GetConnection().GetTCPConnection().Write(data)
		if err != nil {
			logrus.Errorf("[%s] ping failed err:%v", fun, err)
		}
	*/
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fun := "PingRouter.Handle"
	logrus.Infof("[%s] start.....msgId:%d msg:%s\n", fun, request.GetMsgId(), request.GetData())

	data := []byte("ping ping....")
	err := request.GetConnection().Send(1, data)
	if err != nil {
		logrus.Errorf("[%s] ping failed err:%v", fun, err)
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fun := "PingRouter.PostHandle"
	logrus.Infof("[%s] start.....", fun)

	/*
		data := []byte("after ping....")
		_, err := request.GetConnection().GetTCPConnection().Write(data)
		if err != nil {
			logrus.Errorf("[%s] ping failed err:%v", fun, err)
		}
	*/
}

func main() {
	s := znet.NewServer("[zinx v04]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}

/* vim: set tabstop=4 set shiftwidth=4 */
