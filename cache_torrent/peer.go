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

func (p *Peer) RunRecv(m interface {}, connAdd torrent.ConnAdder) {
  /* Backward compatible. */
  p.Peer.RunRecv(m, connAdd)

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

func (p *Peer) AddConnector(id string) {
  connector := torrent.
    NewConnector(p.Id, id, p.Components).
    WithHandshake().
    WithUpload().
    WithDownload()

  p.Connectors[id] = connector
  p.Manager.AddConnector(connector)

  go connector.Run()
}

func (p *Peer) AddLeaderLocalConnector(id string) {
  connector := torrent.
    NewConnector(p.Id, id, p.Components).
    WithHandshake().
    WithUpload().
    WithDownload()

  p.IndirectConnectors[id] = connector
  go connector.Run()
}

func (p *Peer) AddLeaderPeerConnector(id string) {
  connector := torrent.
    NewConnector(p.Id, id, p.Components).
    WithHandshake().
    WithUpload().
    WithDownload()

  p.IndirectConnectors[id] = connector
  go connector.Run()
}
