package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
  "github.com/danalex97/nfsTorrent/log"
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

  proxy.Peer      = proxy.Peer.New(util).(*cache_torrent.Peer)

  proxy.Id        = proxy.Id + ":" + id
  proxy.Transport = NewTransportProxy(proxy.Transport)

  return proxy
}

func (p *PeerProxy) Init(trackerId string) {
  p.Tracker = trackerId

  log.Println("MultiTorrent node", p.Id, "started with tracker", p.Tracker)

  // p.Transport.ControlSend(p.Tracker, Join{p.Id})
}
