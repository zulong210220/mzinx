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
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fun := "PingRouter.Handle"
	logrus.Infof("[%s] start.....msgId:%d msg:%s %d\n", fun, request.GetMsgId(), request.GetData(), request.GetConnection().GetConnID())

	data := []byte("ping ping....")
	err := request.GetConnection().Send(1, data)
	if err != nil {
		logrus.Errorf("[%s] ping failed err:%v", fun, err)
	}
}

type HelloRouter struct {
}

func (pr *HelloRouter) PostHandle(request ziface.IRequest) {
	fun := "HelloRouter.PostHandle"
	logrus.Infof("[%s] start.....", fun)

}

func (pr *HelloRouter) PreHandle(request ziface.IRequest) {
	fun := "HelloRouter.PreHandle"
	logrus.Infof("[%s] start.....msgId:%d msg:%s", fun, request.GetMsgId(), request.GetData())
}

func (pr *HelloRouter) Handle(request ziface.IRequest) {
	fun := "HelloRouter.Handle"
	logrus.Infof("[%s] start.....msgId:%d msg:%s %d\n", fun, request.GetMsgId(), request.GetData(), request.GetConnection().GetConnID())

	data := []byte("hello hello....")
	err := request.GetConnection().Send(1, data)
	if err != nil {
		logrus.Errorf("[%s] ping failed err:%v", fun, err)
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fun := "PingRouter.PostHandle"
	logrus.Infof("[%s] start.....", fun)

}

func DoConnStart(conn ziface.IConnection) {
	fun := "DoConnStart"
	logrus.Errorf("[%s] called ....", fun)

	conn.SetProperty("name", "property test")

	err := conn.Send(2, []byte("DoConnStart BEGIN......."))
	logrus.Errorf("[%s] called ....ok", fun)
	if err != nil {
		logrus.Errorf("[%s] Send failed err:%v", fun, err)
	}
}

func DoConnStop(conn ziface.IConnection) {
	fun := "DoConnStop"
	logrus.Errorf("[%s] called ....", fun)

	name, err := conn.GetProperty("name")
	logrus.Infof("[%s] GetProperty value:%v err:%v", fun, name, err)
}

func main() {
	s := znet.NewServer("[zinx v04]")

	s.SetOnConnStart(DoConnStart)
	s.SetOnConnStop(DoConnStop)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}

/* vim: set tabstop=4 set shiftwidth=4 */
