package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "github.com/danalex97/nfsTorrent/log"
)

// The CacheTorrent peer is only a container extending the BitTorrent peer.
// Depending on the role assigned by the election component, the CacheTorrent
// peer will either implement a Follower or a Leader. The Follower mostly acts
// like a BitTorrent peer, whereas the Leader will have further modifications
// such as message forwarding and special piece picking algorithm, basically
// acting as a proxy between followers and other peers.
type Peer struct {
  *torrent.Peer

  Leaders  []string
  amLeader bool

  node   torrent.Runner
}

func (p *Peer) New(util TorrentNodeUtil) TorrentNode {
  peer := new(Peer)
  peer.Peer     = (peer.Peer.New(util)).(*torrent.Peer)
  peer.amLeader = false
  return peer
}

func (p *Peer) OnJoin() {
  if p.Transport == nil {
    return
  }

  go func() {
    p.Init()
    go p.CheckMessages(p.Bind, p.Process)
  }()
}

// The Process method is identical to BitTorrent's one, but it dispaches the
// message towards the specific node component(i.e. Leader of Follower).
func (p *Peer) Process(m interface {}, state int) {
  switch state {
  case torrent.BindRun:
    if p.amLeader {
      p.node = NewLeader(p)
    } else {
      p.node = NewFollower(p)
    }
    p.node.Run()
  case torrent.BindRecv:
    p.node.Recv(m)
  }
}

// The Bind method is backwards compatible with BitTorrent's Bind, but it
// support 'cache_torrent.Neighbours' messages and 'cache_torrent.Leaders'
// messages.
func (p *Peer) Bind(m interface {}) (state int) {
  switch msg := m.(type) {
  // -- New Protocol --
  case Neighbours:
    // Location awareness extension
    state = torrent.BindDone
    p.Ids = msg.Ids

    // Candidate in the Leader Election
    p.Transport.ControlSend(p.Tracker, Candidate{
      Id   : p.Id,
      Up   : p.Transport.Up(),
      Down : p.Transport.Down(),
    })
  case Leaders:
    state = torrent.BindDone
    p.Leaders = msg.Ids
    log.Println(p.Id, "has Leaders", p.Leaders)

    if isLeader(p) {
      p.amLeader = true
    }
    // Check if I am a seed
    p.Transport.ControlSend(p.Tracker, torrent.SeedReq{p.Id})
  // -- Backward compatible --
  default:
    state = p.Peer.Bind(m)
  }
  return
}

func (p *Peer) GetId(m interface {}) string {
  switch msg := m.(type) {
  case LeaderStart:
    return msg.Id
  default:
    return p.Peer.GetId(m)
  }
}
