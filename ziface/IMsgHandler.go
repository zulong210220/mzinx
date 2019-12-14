package ziface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgId uint32, router IRouter)
}
