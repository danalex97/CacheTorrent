package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "math/rand"
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
  go p.CheckMessages(p.Bind, p.Process)
}

func (p *Peer) Process(m interface {}, state int) {
  switch state {
  case torrent.BindRun:
    p.Run(p.AddConnector)
  case torrent.BindRecv:
    p.RunRecv(m, p.AddConnector)
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

    // Check if I am a seed
    p.Transport.ControlSend(p.Tracker, torrent.SeedReq{p.Id})
  /* Backward compatible. */
  default:
    state = p.Peer.Bind(m)
  }
  return
}

func (p *Peer) amLeader() bool {
  for _, id := range p.Leaders {
    if id == p.Id {
      return true
    }
  }
  return false
}

func (p *Peer) AddConnector(id string) {
  if getAS(p.Id) == getAS(id) || p.amLeader() {
    // Connection within the same AS
    p.Peer.AddConnector(id)
  } else {
    // Connection in different AS

    leader := p.Leaders[rand.Intn(len(p.Leaders))]
    connector := Extend(torrent.
      NewConnector(p.Id, leader, p.Components).
      WithHandshake().
      WithUpload()).
      WithIndirectDownload()

    // Wire the connector
    p.Manager.AddConnector(connector.Connector)
    // We register the connection for the distant peer, so
    // we need to overwrite the sender at the Border node
    p.Connectors[id] = connector

    go connector.Run()

    // Start Inidirect Connection
    p.Transport.ControlSend(leader, LeaderStart{
      Id   : p.Id,
      Dest : id,
    })
  }
}
