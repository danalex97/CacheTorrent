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
    t.ControlSend(conn.to, Have{conn.from, index})
  }
}

func (m *ConnectionManager) Uploads() (uploads []Upload) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    uploads = append(uploads, conn.upload)
  }
  return
}

func (m *ConnectionManager) Downloads() (downloads []Download) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    downloads = append(downloads, conn.download)
  }
  return
}
