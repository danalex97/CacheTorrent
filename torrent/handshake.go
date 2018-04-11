package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "runtime"
  "sync"
)

type Handshake struct {
  *Components
  *sync.RWMutex

  from  string
  to    string

  sent  bool
  done  bool
  link  Link
}

func NewHandshake(connector *Connector) *Handshake {
  return &Handshake{
    connector.components,
    new(sync.RWMutex),
    connector.from,
    connector.to,
    false,
    false,
    (Link)(nil),
  }
}

func (h *Handshake) Link() Link {
  h.RLock()
  for !h.done {
    h.RUnlock()
    runtime.Gosched()
  }

  return h.link
}

func (h *Handshake) Run() {
  h.Lock()
  defer h.Unlock()

  if h.link == nil {
    h.link = h.Transport.Connect(h.to)
    h.Transport.ControlSend(h.to, connReq{h.from, h.link})

    h.sent = true
  }
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

  // I receive a request, but I sent a link.
  // It must be a bidirectional connection.
  if h.sent {
    // We agree on using the link of the smallest id.
    if h.from > h.to {
      h.link = req.link
    }
  } else {
    // I did not send a request, thus I will use the other's connection
    h.link = req.link

    // Let the intiator know handshake is done
    h.Transport.ControlSend(h.to, connRes{h.from})
  }

  h.done = true
}

func (h *Handshake) handleRes(req connRes) {
  h.Lock()
  defer h.Unlock()

  h.done = true
}
