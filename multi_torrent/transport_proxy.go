package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "strings"
)

func ExternId(id string) string {
  if !strings.Contains(id, ":") {
    return id
  }
  return strings.Split(id, ":")[0]
}

func FullId(id string, id2 string) string {
  return id + ":" + id2
}

func InternId(id string) string {
  if !strings.Contains(id, ":") {
    return ""
  }
  return strings.Split(id, ":")[1]
}

type StripProxy struct {
  Transport
}

func NewStripProxy(t Transport) *StripProxy {
  return &StripProxy{
    Transport : t,
  }
}

func (t *StripProxy) Connect(id string) Link {
  return t.Transport.Connect(ExternId(id))
}

func (t *StripProxy) ControlPing(id string) bool {
  return t.Transport.ControlPing(ExternId(id))
}

func (t *StripProxy) ControlSend(id string, m interface {}) {
  t.Transport.ControlSend(ExternId(id), m)
}
