package znet

import (
	"mzinx/config"
	"mzinx/ziface"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
	wg             *sync.WaitGroup
	exitChan       chan bool
	//count int
}

func NewMsgHandler() *MsgHandler {
	mh := &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: config.GetConfig().WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, config.GetConfig().WorkerPoolSize),
		wg:             &sync.WaitGroup{},
		exitChan:       make(chan bool),
	}

	if config.GetConfig().IsWorker {
		mh.StartWorkerPool()
	}

	return mh
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

func (mh *MsgHandler) SendMsg2TaskQueue(request ziface.IRequest) {
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	mh.TaskQueue[workerId] <- request
}

func (mh *MsgHandler) Work(workerId int, taskQueue chan ziface.IRequest) {
	fun := "MsgHandler.Work"
	logrus.Infof("[%s] starting...", fun)
	mh.wg.Add(1)
	defer logrus.Infof("[%s] %d ending...", fun, workerId)
	defer mh.wg.Done()
	tm := time.NewTicker(3 * time.Second)

	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
			// 处理完毕之后再退出
		case <-mh.exitChan:
			logrus.Infof("[%s] %d stoping...", fun, workerId)
			goto END
		case <-tm.C:
			//logrus.Infof("[%s] %d running...", fun, workerId)
		}
	}
END:
	for request := range taskQueue {
		mh.DoMsgHandler(request)
	}
	//mh.count++
}

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, config.GetConfig().TaskQueueSize)
		go mh.Work(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) Stop() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.exitChan <- true
	}
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		close(mh.TaskQueue[i])
	}
	mh.wg.Wait()
	close(mh.exitChan)
	// 测试是否全部退出
	//os.Exit(mh.count)
}
