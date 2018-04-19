package cache_torrent

import (
  "fmt"
)

type Leader struct {
  *Peer
}

func NewLeader(p *Peer) *Leader {
  return &Leader{
    Peer : p,
  }
}

func (l *Leader) Run() {
  fmt.Println("Leader running.")
}

func (l *Leader) Recv(m interface {}) {
}
