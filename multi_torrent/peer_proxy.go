package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
  "github.com/danalex97/nfsTorrent/torrent"
  "github.com/danalex97/nfsTorrent/log"
)

// A PeerProxy is a Peer wrapper over Peer which is used to expose a more useful
// interface towards MultiPeer. The structure has no specific features besides
// providing more methods. The PeerProxy uses a TransportProxy to relay
// messages.
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

// Method called after MultiPeer initalization.
func (p *PeerProxy) Init(trackerId string) {
  p.Tracker = trackerId

  log.Println("MultiTorrent node", p.Id, "started with tracker", p.Tracker)
}

// Method used to strip the pieces from an incoming message to the pieces that
// only subnode PeerProxy is responsible for.
func (p *PeerProxy) SetPieces(pieces []torrent.PieceMeta) {
  p.Pieces = []torrent.PieceMeta{}
  for _, piece := range pieces {
    if piece.Index >= p.piecesFrom && piece.Index < p.piecesTo {
      p.Pieces = append(p.Pieces, piece)
    }
  }
}

// Method used to decorate piece Ids with the <multipeer-id>.
func (p *PeerProxy) SetIds(ids []string) {
  for _, id := range ids {
    p.Ids = append(p.Ids, FullId(id, p.id))
  }
}

// Method used to use the BindRun callback only once.
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
