package znet

import (
	"errors"
	"mzinx/ziface"
	"sync"

	"github.com/sirupsen/logrus"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.Mutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	fun := "ConnManager.Add"
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	connId := conn.GetConnID()
	cm.connections[connId] = conn
	logrus.Infof("[%s] conn:%d ok len:%d", fun, connId, cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	fun := "ConnManager.Remove"
	connId := conn.GetConnID()
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, connId)
	logrus.Infof("[%s] connId:%d success len:%d", fun, connId, cm.Len())
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	fun := "ConnManager.ClearConn"
	tmpConns := make(map[uint32]ziface.IConnection)

	cm.connLock.Lock()
	for connId, conn := range cm.connections {
		tmpConns[connId] = conn
	}
	cm.connLock.Unlock()

	// 解决先关闭server死锁的问题
	for connId, conn := range tmpConns {
		conn.Stop()
		delete(cm.connections, connId)
	}

	logrus.Infof("[%s] success len:%d", fun, cm.Len())
}
func (cm *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	// fun := "ConnManager.Get"
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	conn, ok := cm.connections[connId]
	if ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}
