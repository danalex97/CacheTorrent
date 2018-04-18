package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "fmt"
)

type Peer struct {
  *torrent.Peer

  Local []string
}

func (p *Peer) New(util TorrentNodeUtil) TorrentNode {
  peer := new(Peer)
  peer.Peer = (peer.Peer.New(util)).(*torrent.Peer)
  return peer
}

func (p *Peer) OnJoin() {
  if p.Transport == nil {
    return
  }

  p.Init()
  go p.CheckMessages(p.Bind)
}

func (p *Peer) Bind(m interface {}) (any bool) {
  switch msg := m.(type) {
  /* Backward compatible. */
  case torrent.TrackerReq:
    any = p.Peer.Bind(m)
  case torrent.Neighbours:
    any = p.Peer.Bind(m)
  case torrent.SeedRes:
    any = p.Peer.Bind(m)
  /* New Protocol. */
  case Neighbours:
    // Location awareness extension
    any = true
    p.Ids   = msg.Ids
    p.Local = msg.Local

    fmt.Println("Local:", p.Local)

    // Find if I'm a seed
    p.Transport.ControlSend(p.Tracker, torrent.SeedReq{p.Id})
  default:
    any = p.Peer.Bind(m)
  }
  return
}

func (p *Peer) RunRecv(m interface {}) {
  /* Backward compatible. */
  p.Peer.RunRecv(m)

  /* New Protocol. */
  // [TODO]
}
