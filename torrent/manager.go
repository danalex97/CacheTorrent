package torrent

import (
  "sync"
)

type Manager struct {
  sync.RWMutex

  conns     []*Connector
}

func NewManager() *Manager {
  return &Manager{
    conns : []*Connector{},
  }
}

func (m *Manager) AddConnector(conn *Connector) {
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

func (m *Manager) Uploads() (uploads []*Upload) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    uploads = append(uploads, conn.upload)
  }
  return
}

func (m *Manager) Downloads() (downloads []*Download) {
  m.RLock()
  defer m.RUnlock()

  for _, conn := range m.conns {
    downloads = append(downloads, conn.download)
  }
  return
}
