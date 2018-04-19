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
  switch msg := m.(type) {
  case LeaderStart:
    p.AddLeaderConnector(msg.Id, msg.Dest)
  case RemoteStart:
    connAdd(msg.Id)
  default:
    p.Peer.RunRecv(m, connAdd)
  }
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
  if getAS(p.Id) == getAS(id) {
    // Connection within the same AS
    p.Peer.AddConnector(id)
  } else {
    // Connection in different AS
    leader := p.Leaders[rand.Intn(len(p.Leaders))]
    connector := torrent.
      NewConnector(p.Id, leader, p.Components).
      WithHandshake().
      WithUpload().
      WithDownload()

    // Wire the connector
    p.Manager.AddConnector(connector)
    // We register the connection for the distant peer, so
    // we need to overwrite the sender at the Border node
    p.Connectors[id] = connector

    go connector.Run()

    // Start Indirect Connection with remote Peer id through
    // the Leader leader
    p.Transport.ControlSend(leader, LeaderStart{
      Id   : p.Id,
      Dest : id,
    })
  }
}

func (p *Peer) AddLeaderConnector(local, remote string) {
  /*
   * So far a leader is unable to accept connections.
   */

  // Add connection between Leader and Local peer
  func () {
    connector := Extend(torrent.
      NewConnector(p.Id, local, p.Components).
      WithHandshake().
      WithDownload()).
      WithUploadWithRedirect()
    p.Manager.AddConnector(connector.Strip())
    p.Connectors[local] = connector.Strip()
    go connector.Run()
  }()

  // Add connection between leader and remote
  func () {
    if _, ok := p.Connectors[remote]; !ok {
      connector := Extend(torrent.
        NewConnector(p.Id, remote, p.Components).
        WithHandshake()).
        WithDownloadWithRedirect()
      p.Connectors[remote] = connector.Strip()

      p.Manager.AddConnector(connector.Strip())
      go connector.Run()

      // Start the Remote connection
      p.Transport.ControlSend(remote, RemoteStart{p.Id})
    }

    connector := p.Connectors[remote].(*torrent.Connector)

    /* Add redirects. */
    switch download := connector.Download.(type) {
    case *download:
      download.AddRedirect(local)
    default:
      connector.Download = NewDownloadWithRedirect(Extend(connector)).
        AddRedirect(local)
    }
  }()
}
