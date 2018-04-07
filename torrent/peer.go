package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "fmt"
)

type Peer struct {
}

func (p *Peer) OnJoin() {
  fmt.Println("Peer")
}

func (p *Peer) OnLeave() {
}

func (p *Peer) New(TorrentNodeUtil) TorrentNode {
  return new(Peer)
}
