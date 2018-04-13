package torrent

import (
  "sync"
)

type Manager struct {
  *sync.Mutex

  conns     []*Connector
}

func NewManager() *Manager {
  return &Manager{
    Mutex : new(sync.Mutex),
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
  m.Lock()
  defer m.Unlock()

  for _, conn := range m.conns {
    uploads = append(uploads, conn.upload)
  }
  return
}

/**
 * We moved some of the responsibility in 'MultiDownload.py',
 * 'download.py' and 'RequestManager.py' in the Manager as we only
 * need a struct which references the list of connections.
 */
func (m *Manager) Lost() {
  m.Lock()
  defer m.Unlock()

  for _, conn := range m.conns {
    // We try to request more pieces only if the connection is not choked
    if !conn.download.choked {
      conn.RequestMore()
    }
  }
}

func (m *Manager) Have(index int) {
  m.Lock()
  defer m.Unlock()

  for _, conn := range m.conns {
    t := conn.components.Transport
    t.ControlSend(conn.to, have{conn.from, index})
  }
}
