package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
  "github.com/danalex97/nfsTorrent/torrent"
  "github.com/danalex97/nfsTorrent/config"

  "strconv"
  "reflect"
  "fmt"
)

var pieceNumber config.Const = config.NewConst(config.Pieces)
var multi       config.Const = config.NewConst(config.Multi)

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
  pieceNbr      := totalPieceNbr / multi.Int()
  for i := 0; i < multi.Int(); i++ {
    piecesFrom := pieceNbr * i
    piecesTo   := piecesFrom + pieceNbr
    if totalPieceNbr < piecesTo {
      piecesTo = totalPieceNbr
    }

    internalId := strconv.Itoa(i)

    // Register new peer proxy
    peer := NewPeerProxy(p.util, internalId, piecesFrom, piecesTo)
    p.peers[internalId] = peer
  }

  go func() {
    // We initalize the Tracker requests only once
    p.Init()

    // Send Join request
    p.Transport.ControlSend(p.Tracker, Join{p.Id})

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

func (p *MultiPeer) getId(m interface {}) string {
  switch msg := m.(type) {
  case cache_torrent.Leaders:
    if len(msg.Ids) == 0 {
      return ""
    } else {
      return msg.Ids[0]
    }
  default:
    return p.Peer.GetId(m)
  }
}

// The MultiTorrent's Bind method distributes all the initialization messages
// towards the PeerProxies.
func (p *MultiPeer) Bind(m interface {}) int {
  id := InternId(p.getId(m))
  if peer, ok := p.peers[id]; ok {
    return peer.Bind(m)
  } else {
    switch m.(type) {
    case torrent.TrackerReq:
      return p.Peer.Bind(m)

    case torrent.SeedRes:
      ret := p.Peer.Bind(m)
      for _, peer := range p.peers {
        peer.SetPieces(p.Pieces)
      }
      return ret

    case cache_torrent.Neighbours:
      ret := p.Peer.Bind(m)
      for _, peer := range p.peers {
        peer.SetIds(p.Ids)
      }
      return ret

    default:
      fmt.Println("Unexpected messsage.", m, reflect.TypeOf(m).String())
      return 0
    }
  }
  return 0
}

// The MultiPeer's Process method redirects all the messages towards the
// correct PeerProxy.
func (p *MultiPeer) Process(m interface {}, state int) {
  id := InternId(p.getId(m))
  if peer, ok := p.peers[id]; ok {
    peer.Process(m, state)
  } else {
    // Broadcast the message if it is not adressed to
    // a particular node
    for _, peer := range p.peers {
      peer.Process(m, state)
    }
  }
}
