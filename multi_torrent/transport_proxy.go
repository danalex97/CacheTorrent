package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "strings"
)

type TransportProxy struct {
  Transport
}

func NewTransportProxy(t Transport) *TransportProxy {
  return &TransportProxy{
    Transport : t,
  }
}

func ExternId(id string) string {
  if !strings.Contains(id, ":") {
    return id
  }
  return strings.Split(id, ":")[0]
}

func FullId(id string, id2 string) string {
  return id + ":" + id2
}

func (t *TransportProxy) Connect(id string) Link {
  return t.Transport.Connect(ExternId(id))
}

func (t *TransportProxy) ControlPing(id string) bool {
  return t.Transport.ControlPing(ExternId(id))
}

func (t *TransportProxy) ControlSend(id string, m interface {}) {
  t.Transport.ControlSend(ExternId(id), m)
}
