package cache_torrent

import (
  "fmt"
)

type Follower struct {
  *Peer
}

func NewFollower(p *Peer) *Follower {
  return &Follower{
    Peer : p,
  }
}

func (l *Follower) Run() {
  fmt.Println("Follower running.")
}

func (l *Follower) Recv(m interface {}) {
}
