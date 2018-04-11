package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "sync"
)

type Handshake struct {
  *Components
  *sync.Mutex

  from  string
  to    string

  sent  bool
  link  Link
}

func NewHandshake(connector *Connector) *Handshake {
  return &Handshake{
    connector.components,
    new(sync.Mutex),
    connector.from,
    connector.to,
    false,
    (Link)(nil),
  }
}

func (h *Handshake) Link() Link {
  return (Link)(nil)
}

func (h *Handshake) Run() {
  h.Lock()
  defer h.Unlock()

}

func (h *Handshake) Recv(m interface {}) {
  switch msg := m.(type) {
  case connReq:
    h.handleReq(msg)
  case connRes:
    h.handleRes(msg)
  }
}

func (h *Handshake) handleReq(req connReq) {
  h.Lock()
  defer h.Unlock()

}

func (h *Handshake) handleRes(res connRes) {
  h.Lock()
  defer h.Unlock()

}
