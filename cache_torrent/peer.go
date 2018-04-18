package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
)

type Peer struct {
  *torrent.Peer
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
  switch m.(type) {
  case torrent.TrackerReq:
    p.Peer.Bind(m)
  case torrent.Neighbours:
    p.Peer.Bind(m)
  case torrent.SeedRes:
    p.Peer.Bind(m)
  default:
    p.Peer.Bind(m)
  }
  return
}

func (p *Peer) RunRecv(m interface {}) {
  p.Peer.RunRecv(m)
}
