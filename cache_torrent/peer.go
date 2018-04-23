package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "fmt"
)

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

  p.Init()
  go p.CheckMessages(p.Bind, p.Process)
}

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

func (p *Peer) Bind(m interface {}) (state int) {
  switch msg := m.(type) {
  /* New Protocol. */
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
    fmt.Println(p.Id, "has Leaders", p.Leaders)

    if isLeader(p) {
      p.amLeader = true
    }
    // Check if I am a seed
    p.Transport.ControlSend(p.Tracker, torrent.SeedReq{p.Id})
  /* Backward compatible. */
  default:
    state = p.Peer.Bind(m)
  }
  return
}

func (p *Peer) GetId(m interface {}) string {
  switch msg := m.(type) {
  case LeaderStart:
    return msg.Id
  case Miss:
    return msg.Id
  default:
    return p.Peer.GetId(m)
  }
}
