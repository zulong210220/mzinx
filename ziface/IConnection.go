/*
 * File : connection.go
 * CreateDate : 2019-12-07 16:35:51
 * */
package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() string
	Send(id uint32, data []byte) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error

/* vim: set tabstop=4 set shiftwidth=4 */
