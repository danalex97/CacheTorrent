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

  // BitTorrent protocol
  pieces []pieceMeta

  // BitTorrent components
  connectors map[string]Runner

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
  peer.connectors = make(map[string]Runner)

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

        // make connectors
        for _, id := range p.ids {
          p.connectors[id] = NewConnector(p.id, id)
        }
      case seedRes:
        p.pieces = msg.pieces

        go p.Run()
      default:
        // We run this only when the connectors are initialized
        if len(p.connectors) > 0 {
          p.RunRecv(m)
        }
      }

    default:
      // allow other nodes in simulation run
      runtime.Gosched()
    }
  }
}

func (p *Peer) Run() {
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
}

func (p *Peer) RunRecv(m interface {}) {
  var connector Runner

  // Redirect the message to the connector
  switch msg := m.(type) {
  case choke: connector = p.connectors[msg.id]
  case unchoke: connector = p.connectors[msg.id]
  case interested: connector = p.connectors[msg.id]
  case notInterested: connector = p.connectors[msg.id]
  case have: connector = p.connectors[msg.id]
  case request: connector = p.connectors[msg.id]
  case piece: connector = p.connectors[msg.id]
  }

  if connector != nil {
    connector.Recv(m)
  }
}
