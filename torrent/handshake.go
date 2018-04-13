package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "runtime"
  "sync"
)

type Handshake struct {
  sync.RWMutex
  *Components

  from  string
  to    string

  done  bool

  downlink Link
  uplink   Link
}

func NewHandshake(connector *Connector) *Handshake {
  t    := connector.components.Transport
  link := t.Connect(connector.to)

  return &Handshake{
    Components: connector.components,
    from:       connector.from,
    to:         connector.to,
    done:       false,
    downlink:   (Link)(nil),
    uplink:     link,
  }
}

func (h *Handshake) Uplink() Link {
  return h.uplink
}

func (h *Handshake) Downlink() Link {
  h.wait()

  h.RLock()
  defer h.RUnlock()

  return h.downlink
}

func (h *Handshake) wait() {
  h.RLock()
  for !h.done {
    h.RUnlock()
    runtime.Gosched()
    h.RLock()
  }
  h.RUnlock()
}

func (h *Handshake) Run() {
  h.Transport.ControlSend(h.to, connReq{h.from, h.uplink})
}

func (h *Handshake) Recv(m interface {}) {
  switch msg := m.(type) {
  case connReq:
    h.handleReq(msg)
  }
}

func (h *Handshake) handleReq(req connReq) {
  h.Lock()
  defer h.Unlock()

  h.downlink = req.link
  h.done = true
}

func (h *Handshake) Done() bool {
  h.RLock()
  defer h.RUnlock()

  return h.done
}
