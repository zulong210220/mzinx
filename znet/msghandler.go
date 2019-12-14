package znet

import (
	"mzinx/ziface"
	"strconv"

	"github.com/sirupsen/logrus"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	fun := "MsgHandler.AddRouter"
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api msgId:%d" + strconv.Itoa(int(msgId)))
	}

	mh.Apis[msgId] = router
	logrus.Infof("[%s] msgId:%d", fun, msgId)
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	fun := "MsgHandler.DoMsgHandler"
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		logrus.Errorf("[%s] api msgId:%d not found", fun, request.GetMsgId())
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
