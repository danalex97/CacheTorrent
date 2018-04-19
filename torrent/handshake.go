package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "runtime"
  "sync"
)

type Handshake interface {
  Runner

  Uplink()   Link
  Downlink() Link
  Done()     bool
}

type handshake struct {
  sync.RWMutex
  *Components

  from  string
  to    string

  done  bool

  downlink Link
  uplink   Link
}

func NewHandshake(connector *Connector) Handshake {
  t    := connector.Transport
  link := t.Connect(connector.To)

  return &handshake{
    Components: connector.Components,
    from:       connector.From,
    to:         connector.To,
    done:       false,
    downlink:   (Link)(nil),
    uplink:     link,
  }
}

func (h *handshake) Uplink() Link {
  return h.uplink
}

func (h *handshake) Downlink() Link {
  h.wait()

  h.RLock()
  defer h.RUnlock()

  return h.downlink
}

func (h *handshake) wait() {
  h.RLock()
  for !h.done {
    h.RUnlock()
    runtime.Gosched()
    h.RLock()
  }
  h.RUnlock()
}

func (h *handshake) Run() {
  h.Transport.ControlSend(h.to, ConnReq{h.from, h.uplink})
}

func (h *handshake) Recv(m interface {}) {
  switch msg := m.(type) {
  case ConnReq:
    h.handleReq(msg)
  }
}

func (h *handshake) handleReq(req ConnReq) {
  h.Lock()
  defer h.Unlock()

  h.downlink = req.Link
  h.done = true
}

func (h *handshake) Done() bool {
  h.RLock()
  defer h.RUnlock()

  return h.done
}
