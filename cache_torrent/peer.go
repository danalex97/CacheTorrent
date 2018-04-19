package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "fmt"
)

type Peer struct {
  *torrent.Peer

  Leaders []string

  IndirectConnectors map[string]torrent.Runner
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
  switch msg := m.(type) {
  /* Backward compatible. */
  case torrent.TrackerReq:
    any = p.Peer.Bind(m)
  case torrent.Neighbours:
    any = p.Peer.Bind(m)
  /* New Protocol. */
  case torrent.SeedRes:
    any = true
    p.Pieces = msg.Pieces

    // Since we do Run here, it must be that it will not hang
    p.Run()
  case Neighbours:
    // Location awareness extension
    any = true
    p.Ids = msg.Ids

    // Candidate in the Leader Election
    p.Transport.ControlSend(p.Tracker, Candidate{
      Id   : p.Id,
      Up   : p.Transport.Up(),
      Down : p.Transport.Down(),
    })
  case Leaders:
    p.Leaders = msg.Ids
    fmt.Println(p.Id, "has Leaders", p.Leaders)

    // Check if I am a seed
    p.Transport.ControlSend(p.Tracker, torrent.SeedReq{p.Id})
  default:
    if len(p.Connectors) > 0 {
      any = true

      // All initialized
      p.RunRecv(m)
    } else {
      // Send message to myself
      p.Transport.ControlSend(p.Id, m)
    }
  }
  return
}

func (p *Peer) RunRecv(m interface {}) {
  /* Backward compatible. */
  p.Peer.RunRecv(m)

  /* New Protocol. */
  switch msg := m.(type) {
  case IndirectReq:
    from := msg.From
    dest := msg.Dest

    // Start connection from Leader to Peer.
    if _, ok := p.IndirectConnectors[dest]; !ok {
      p.AddLeaderPeerConnector(dest)
    }

    // Start connection from Leader to Local.
    if _, ok := p.IndirectConnectors[from]; !ok {
      p.AddLeaderLocalConnector(from)
    }

    // Forward request to Leader-Peer connector
    if connector, ok := p.IndirectConnectors[dest]; ok {
      connector.Recv(m)
    }
  }
}

func (p *Peer) AddLeaderLocalConnector(id string) {
  connector := NewLeaderLocalConnector(p.Id, id, p.Components)
  p.IndirectConnectors[id] = connector
  go connector.Run()
}

func (p *Peer) AddLeaderPeerConnector(id string) {
  connector := NewLeaderPeerConnector(p.Id, id, p.Components)
  p.IndirectConnectors[id] = connector
  go connector.Run()
}
