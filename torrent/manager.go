package torrent

import (
  "sync"
)

// The Connection Manager is a lock-guarded list of Connectors. Since
// we assume reliable delivery, unlike in version 5.3, the Connection
// Manager does not have further responsibilities such as keeping connections
// alive.
type Manager interface {
  AddConnector(conn *Connector)

  Uploads() []Upload
  Downloads() []Download
}

type connectionManager struct {
  sync.RWMutex

  conns     []*Connector
}

func NewConnectionManager() Manager {
  return &connectionManager{
    conns : []*Connector{},
  }
}

func (m *connectionManager) AddConnector(conn *Connector) {
  m.Lock()
  defer m.Unlock()

  m.conns = append(m.conns, conn)

  // Send haves at connection
  s := conn.Storage
  t := conn.Transport
  for _, index := range s.Pieces() {
    t.ControlSend(conn.To, Have{conn.From, index})
  }
}

func (m *connectionManager) Uploads() (uploads []Upload) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    if conn.Upload != nil {
      uploads = append(uploads, conn.Upload)
    }
  }
  return
}

func (m *connectionManager) Downloads() (downloads []Download) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    if conn.Download != nil {
      downloads = append(downloads, conn.Download)
    }
  }
  return
}
