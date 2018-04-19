package torrent

import (
  "sync"
)

type Manager interface {
  AddConnector(conn *Connector)
  Uploads() []Upload
  Downloads() []Download
}

type ConnectionManager struct {
  sync.RWMutex

  conns     []*Connector
}

func NewConnectionManager() Manager {
  return &ConnectionManager{
    conns : []*Connector{},
  }
}

func (m *ConnectionManager) AddConnector(conn *Connector) {
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

func (m *ConnectionManager) Uploads() (uploads []Upload) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    if conn.Upload != nil {
      uploads = append(uploads, conn.Upload)
    }
  }
  return
}

func (m *ConnectionManager) Downloads() (downloads []Download) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    if conn.Download != nil {
      downloads = append(downloads, conn.Download)
    }
  }
  return
}
