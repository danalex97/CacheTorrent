package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
)

type PeerProxy struct {
  *cache_torrent.Peer

  id string

  piecesFrom int
  piecesTo   int
}

func NewPeerProxy(util TorrentNodeUtil, id string, piecesFrom, piecesTo int) *PeerProxy {
  proxy := &PeerProxy{
    id : id,

    piecesFrom : piecesFrom,
    piecesTo   : piecesTo,
  }

  return proxy
}

func (p *PeerProxy) FullId() string {
  return p.Id + "." + p.id
}
