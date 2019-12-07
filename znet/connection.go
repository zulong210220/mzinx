/*
 * File : connection.go
 * CreateDate : 2019-12-07 17:01:46
 * */
package znet

import (
	"io"
	"mzinx/consts"
	"mzinx/ziface"
	"net"

	"strings"

	"github.com/sirupsen/logrus"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	//handleApi ziface.HandleFunc
	ExitChan chan bool
	Router   ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connId uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connId,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fun := "Connection.StartReader"
	logrus.Infof("[%s] Reader Groutines is running...\n", fun)

	defer logrus.Infof("[%s] connID:%d stoping", fun, c.ConnID)
	defer c.Stop()

	for {
		buf := make([]byte, consts.BufLineSize)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			logrus.Errorf("[%s] read buffer cnt:%d failed err:%v", fun, cnt, err)

			if err == io.EOF || strings.Contains(err.Error(), "tcp") {
				break
			}
			continue
		}

		/*
			if err := c.handleApi(c.Conn, buf, cnt); err != nil {
				logrus.Errorf("[%s] connID:%d handleApi failed err:%v", fun, c.ConnID, err)
				break
			}
		*/

		req := &Request{
			conn: c,
			data: buf,
		}
		go func(req ziface.IRequest) {
			logrus.Infof("[%s] go func....", fun)
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(req)
	}
}

func (c *Connection) Start() {
	fun := "Connection.Start"
	logrus.Infof("[%s] Conn Start.. ConnID:%d", fun, c.ConnID)
	go c.StartReader()
}

func (c *Connection) Stop() {
	fun := "Connection.Stop"
	logrus.Infof("[%s] Conn Stop.. ConnID:%d", fun, c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() string {
	return c.Conn.RemoteAddr().String()
}

func (c *Connection) Send(data []byte) error {
	return nil
}

/* vim: set tabstop=4 set shiftwidth=4 */
