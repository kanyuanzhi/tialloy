package tiface

type IConnManager interface {
	Add(connection IConnection) error
	Remove(connection IConnection)
	Get(connID string) (IConnection, error)
	Len() int
	ClearAllConn()
}
