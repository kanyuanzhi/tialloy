package tinet

import (
	"errors"
	"fmt"
	"github.com/kanyuanzhi/tialloy/tiface"
	"github.com/kanyuanzhi/tialloy/utils"
	"sync"
)

type ConnManager struct {
	connections map[uint32]tiface.IConnection
	connLock    *sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]tiface.IConnection),
		connLock:    new(sync.RWMutex),
	}
}

func (cm *ConnManager) Add(conn tiface.IConnection) error {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	if _, ok := cm.connections[conn.GetConnID()]; ok {
		return errors.New(fmt.Sprintf("connID=%d is exsited", conn.GetConnID()))
	}
	cm.connections[conn.GetConnID()] = conn
	utils.GlobalLog.Infof("add connID=%d to connManager successfully, current conn num=%d", conn.GetConnID(), cm.Len())
	return nil
}

func (cm *ConnManager) Remove(conn tiface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
	utils.GlobalLog.Warnf("remove connID=%d from connManager successfully, current conn num=%d", conn.GetConnID(), cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (tiface.IConnection, error) {
	cm.connLock.RLocker()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New(fmt.Sprintf("connID=%d is not exsited", connID))
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearAllConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, conn := range cm.connections {
		conn.Stop() // 删除之前先关闭连接
		// TODO:此处应通知客户端服务器关闭连接？
		delete(cm.connections, connID)
	}
	utils.GlobalLog.Tracef("clear all connections from connManager successfully, current conn num=%d", cm.Len())
}
