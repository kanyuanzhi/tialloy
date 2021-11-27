package tiface

type IConnManager interface {
	Add(connection IConnection) error
	Remove(connection IConnection)
	Get(connID uint32) (IConnection, error)
	Len() int
	ClearAllConn()
}


