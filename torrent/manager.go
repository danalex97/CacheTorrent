package torrent

import (
  "sync"
)

type Manager interface {
  AddConnector(conn *Connector)
  Uploads() []*Upload
  Downloads() []*Download
}

type ConnectionManager struct {
  sync.RWMutex

  conns     []*Connector
}

func NewConnectioManager() Manager {
  return &ConnectionManager{
    conns : []*Connector{},
  }
}

func (m *ConnectionManager) AddConnector(conn *Connector) {
  m.Lock()
  defer m.Unlock()

  m.conns = append(m.conns, conn)

  // Send haves at connection
  s := conn.components.Storage
  t := conn.components.Transport
  for _, piece := range s.pieces {
    t.ControlSend(conn.to, have{conn.from, piece.index})
  }
}

func (m *ConnectionManager) Uploads() (uploads []*Upload) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    uploads = append(uploads, conn.upload)
  }
  return
}

func (m *ConnectionManager) Downloads() (downloads []*Download) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    downloads = append(downloads, conn.download)
  }
  return
}
