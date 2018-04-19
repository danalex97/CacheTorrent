package cache_torrent

import (
  "fmt"
)

type Forwarder struct {
  from string
  to   string
}

func NewForwarder(from, to string) *Forwarder {
  return &Forwarder{
    from : from,
    to   : to,
  }
}

func (f *Forwarder) Recv(m interface {}) {
  fmt.Println("Fwd", f.from, f.to, m)
}
