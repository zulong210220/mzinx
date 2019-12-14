/*
 * File : connection.go
 * CreateDate : 2019-12-07 17:01:46
 * */
package znet

import (
	"errors"
	"io"
	"mzinx/ziface"
	"net"

	"encoding/hex"

	"github.com/sirupsen/logrus"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	//handleApi ziface.HandleFunc
	ExitChan chan bool
	msgChan  chan []byte
	handler  ziface.IMsgHandler
}

func NewConnection(conn *net.TCPConn, connId uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connId,
		handler:  handler,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		msgChan:  make(chan []byte),
	}

	return c
}

func (c *Connection) StartReader() {
	fun := "Connection.StartReader"
	logrus.Infof("[%s] Reader Groutines is running...\n", fun)

	defer logrus.Infof("[%s] connID:%d stoping", fun, c.ConnID)
	defer c.Stop()

	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			logrus.Errorf("[%s] read msg head err:%v", fun, err)
			return
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			logrus.Errorf("[%s] Unpack headData:%s err:%v", fun, headData, err)
			return
		}

		logrus.Infof("[%s] ok msg id:%d len:%d", fun, msg.GetMsgId(), msg.GetDataLen())
		data := []byte{}
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				logrus.Errorf("[%s] ReadFull data err:%v", fun, err)
				return
			}
		}
		msg.SetData(data)

		req := &Request{
			conn: c,
			msg:  msg,
		}

		go c.handler.DoMsgHandler(req)
	}
}

func (c *Connection) Start() {
	fun := "Connection.Start"
	logrus.Infof("[%s] Conn Start.. ConnID:%d", fun, c.ConnID)
	go c.StartReader()
	go c.StartWriter()
}

func (c *Connection) Stop() {
	fun := "Connection.Stop"
	logrus.Infof("[%s] Conn Stop.. ConnID:%d", fun, c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	c.Conn.Close()
	close(c.msgChan)
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

func NewMsgPack(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (c *Connection) Send(msgId uint32, data []byte) error {
	fun := "Connection.Send"
	if c.isClosed {
		return errors.New("connection is closed")
	}

	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPack(msgId, data))
	if err != nil {
		logrus.Errorf("[%s] Pack error msgId:%d msg:%s", fun, msgId, data)
		return err
	}

	logrus.Infof("[%s] binary msg:%s :%d Id:%d data:%s", fun, hex.EncodeToString(binaryMsg), len(binaryMsg), msgId, data)

	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) StartWriter() {
	fun := "Connection.StartWriter"
	logrus.Infof("[%s] starting...", fun)
	defer logrus.Infof("[%s] client:%s exit", fun, c.RemoteAddr())

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				logrus.Errorf("[%s] Write data:%s err:%s", fun, data, err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

/* vim: set tabstop=4 set shiftwidth=4 */
