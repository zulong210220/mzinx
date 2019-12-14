/*
 * File : IServer.go
 * CreateDate : 2019-12-07 15:15:26
 * */

package ziface

type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 运行
	Serve()

	AddRouter(msgId uint32, router IRouter)
}

/* vim: set tabstop=4 set shiftwidth=4 */
