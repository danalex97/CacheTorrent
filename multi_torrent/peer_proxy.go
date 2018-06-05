package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
  "github.com/danalex97/nfsTorrent/torrent"
  "github.com/danalex97/nfsTorrent/log"
)

type PeerProxy struct {
  *cache_torrent.Peer

  id string

  piecesFrom int
  piecesTo   int

  bind bool
}

func NewPeerProxy(util TorrentNodeUtil, id string, piecesFrom, piecesTo int) *PeerProxy {
  proxy := &PeerProxy{
    id : id,

    piecesFrom : piecesFrom,
    piecesTo   : piecesTo,

    bind : false,
  }

  proxy.Peer      = proxy.Peer.New(util).(*cache_torrent.Peer)

  proxy.Id        = FullId(proxy.Id, id)
  proxy.Transport = NewStripProxy(proxy.Transport)

  return proxy
}

func (p *PeerProxy) Init(trackerId string) {
  p.Tracker = trackerId

  log.Println("MultiTorrent node", p.Id, "started with tracker", p.Tracker)
}

func (p *PeerProxy) SetPieces(pieces []torrent.PieceMeta) {
  p.Pieces = []torrent.PieceMeta{}
  for _, piece := range pieces {
    if piece.Index >= p.piecesFrom && piece.Index < p.piecesTo {
      p.Pieces = append(p.Pieces, piece)
    }
  }
}

func (p *PeerProxy) SetIds(ids []string) {
  for _, id := range ids {
    p.Ids = append(p.Ids, FullId(id, p.id))
  }
}

func (p *PeerProxy) Process(m interface {}, state int) {
  switch state {
  case torrent.BindRun:
    if !p.bind {
      p.bind = true
      p.Peer.Process(m, state)
    }
  case torrent.BindRecv:
    p.Peer.Process(m, state)
  }
}
