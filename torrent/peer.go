package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "runtime"
  "fmt"
)

type Peer struct {
  id      string
  tracker string
  ids     []string

  transport Transport
  time      func() int

  // BitTorrent protocol
  pieces []pieceMeta

  // BitTorrent components
  connectors    map[string]Runner // the connectors that were chosen by tracker
  allConnectors map[string]Runner // the connections with other peers(uploaders)
  components *Components

  // used only to identify tracker
  join    string
}

/* Implementation of Torrent Node interface. */
func (p *Peer) OnJoin() {
  // If the transport is missing, it must be we are
  // on torrent-less node
  if p.transport == nil {
    return
  }

  p.Init()
  go p.InitRecv()
}

func (p *Peer) OnLeave() {
}

func (p *Peer) New(util TorrentNodeUtil) TorrentNode {
  peer := new(Peer)

  peer.id        = util.Id()
  peer.join      = util.Join()
  peer.transport = util.Transport()

  peer.pieces     = []pieceMeta{}
  peer.connectors    = make(map[string]Runner)
  peer.allConnectors = make(map[string]Runner)
  peer.components = new(Components)

  peer.time = util.Time()

  return peer
}

/* Internal functions. */
func (p *Peer) Init() {
  // Find out who the tracker is
  p.transport.ControlSend(p.join, trackerReq{p.id})

  msg := <-p.transport.ControlRecv()
  p.tracker = msg.(trackerRes).id

  // The peer should be initialized
  fmt.Printf("Node %s started with tracker %s\n", p.id, p.tracker)

  // Send join message to the tracker
  p.transport.ControlSend(p.tracker, join{p.id})
}

func (p *Peer) InitRecv() {
  for {
    select {
    case m, ok := <-p.transport.ControlRecv():
      if !ok {
        continue
      }

      switch msg := m.(type) {
      case trackerReq:
        p.transport.ControlSend(msg.from, trackerRes{p.tracker})
      case neighbours:
        p.ids = msg.ids

        // Find if I'm a seed
        p.transport.ControlSend(p.tracker, seedReq{p.id})
      case seedRes:
        p.pieces = msg.pieces

        // Since we do Run here, it must be that it will not hang
        p.Run()
      default:
        // Wait for connectors to get initialized
        if len(p.connectors) > 0 {
          p.RunRecv(m)
        } else {
          // Requeue the message if the connectors are not initialized
          p.transport.ControlSend(p.id, m)
        }
      }

    default:
      // allow other nodes in simulation run
      runtime.Gosched()
    }
  }
}

func (p *Peer) Run() {
  // We want to bind all these variables here, so
  // we don't need any synchroization.

  // make per peer variables
  p.components.Picker        = NewPicker(p.pieces)
  p.components.Storage       = NewStorage(p.pieces)
  p.components.Transport     = p.transport
  p.components.Choker        = NewChoker(p.time)
  p.components.MultiDownload = NewMultiDownload()

  // make connectors
  for _, id := range p.ids {
    connector := NewConnector(p.id, id, p.components)
    p.components.Choker.AddConnector(connector)

    p.connectors[id]    = connector
    p.allConnectors[id] = p.connectors[id]
  }

  // Let all the neighbouring peers know what pieces I have
  for _, id := range p.ids {
    for _, piece := range p.pieces {
      p.transport.ControlSend(id, have{p.id, piece.index})
    }
  }

  // Run all connectors
  for _, connector := range p.connectors {
    go connector.Run()
  }

  go p.components.Choker.Run()
}

func (p *Peer) RunRecv(m interface {}) {
  id := ""

  // Redirect the message to the connector
  switch msg := m.(type) {
  case choke: id = msg.id
  case unchoke: id = msg.id
  case interested: id = msg.id
  case notInterested: id = msg.id
  case have: id = msg.id
  case request: id = msg.id
  case piece: id = msg.id
  }

  if id == "" {
    return
  }

  if _, ok := p.allConnectors[id]; !ok {
    p.allConnectors[id] = NewConnector(p.id, id, p.components)
  }

  connector := p.allConnectors[id]
  connector.Recv(m)
}
