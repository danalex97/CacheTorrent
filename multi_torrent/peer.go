package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
  "github.com/danalex97/nfsTorrent/config"

  "strconv"
)

var pieceNumber config.Const = config.NewConst(config.Pieces)

var MultiPeerMembers int = 1

// A MultiPeer is wrapper over multiple Peers which follow the CacheTorrent
// protocol. The Original message IDs are decorated with an internal ID, each
// MultiPeer having an ID format "<multipeer-id>.<peer-id>". Once a message
// arrives to a Peer, the id is stripped to "<multipeer-id>".
type MultiPeer struct {
  *cache_torrent.Peer

  // Map from internal ID to Peer.
  peers map[string]*PeerProxy

  // Utility structure to pass at single Peer initilization.
  util TorrentNodeUtil
}

func (p *MultiPeer) New(util TorrentNodeUtil) TorrentNode {
  return &MultiPeer{
    peers : make(map[string]*PeerProxy),
    util  : util,
  }
}

func (p *MultiPeer) OnJoin() {
  if p.Transport == nil {
    return
  }

  totalPieceNbr := pieceNumber.Int()
  pieceNbr      := totalPieceNbr / MultiPeerMembers
  for i := 0; i < MultiPeerMembers; i++ {
    piecesFrom := pieceNbr * i
    piecesTo   := piecesFrom + pieceNbr
    if totalPieceNbr < piecesTo {
      piecesTo = totalPieceNbr
    }

    internalId := strconv.Itoa(i)

    // Register new peer proxy
    peer := NewPeerProxy(p.util, internalId, piecesFrom, piecesTo)
    p.peers[peer.Id] = peer
  }
}

func (p *MultiPeer) OnLeave() {
}
