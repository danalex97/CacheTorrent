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
  multiPeer := &MultiPeer{
    peers : make(map[string]*PeerProxy),
    util  : util,
  }

  multiPeer.Peer = multiPeer.Peer.New(util).(*cache_torrent.Peer)

  return multiPeer
}

func (p *MultiPeer) OnJoin() {
  if p.Transport == nil {
    return
  }

  // Build all the PeerProxies
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

  go func() {
    // We initalize the Tracker requests only once
    p.Init()

    // Since Join messages will be ignored by the new Tracker,
    // we will run the initialization for the PeerProxies
    for _, proxy := range p.peers {
      proxy.Init(p.Tracker)
    }

    // Start checking messages and redirecting them to the
    // respective PeerProxy
    go p.CheckMessages(p.Bind, p.Process)
  }()
}

func (p *MultiPeer) Bind(m interface {}) int {
  id := p.GetId(m)
  return p.peers[id].Bind(m)
}

func (p *MultiPeer) Process(m interface {}, state int) {
  id := p.GetId(m)
  p.peers[id].Process(m, state)
}
