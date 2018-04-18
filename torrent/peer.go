package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "runtime"
  // "reflect"
  "fmt"
)

const maxPeers int = config.InPeers + config.OutPeers

type Peer struct {
  *Components

  id      string
  tracker string
  ids     []string

  time      func() int

  // BitTorrent protocol
  pieces []pieceMeta

  // BitTorrent components
  connectors  map[string]Runner // the connectors that were chosen by tracker

  // used only to identify tracker
  join    string
}

/* Implementation of Torrent Node interface. */
func (p *Peer) OnJoin() {
  // If the Transport is missing, it must be we are
  // on torrent-less node
  if p.Transport == nil {
    return
  }

  p.Init()
  go p.CheckMessages(p.Bind)
}

func (p *Peer) OnLeave() {
}

func (p *Peer) New(util TorrentNodeUtil) TorrentNode {
  peer := new(Peer)

  peer.id        = util.Id()
  peer.join      = util.Join()
  peer.Components = new(Components)
  peer.Transport  = util.Transport()

  peer.pieces     = []pieceMeta{}
  peer.connectors = make(map[string]Runner)

  peer.time = util.Time()

  return peer
}

/* Internal functions. */
func (p *Peer) Init() {
  // Find out who the tracker is
  p.Transport.ControlSend(p.join, TrackerReq{p.id})

  msg := <-p.Transport.ControlRecv()
  p.tracker = msg.(TrackerRes).id

  // The peer should be initialized
  fmt.Printf("Node %s started with tracker %s\n", p.id, p.tracker)

  // Send join message to the tracker
  p.Transport.ControlSend(p.tracker, Join{p.id})
}

func (p *Peer) CheckMessages(process func(interface {}) bool) {
  for {
    messages := []interface{}{}

    // Get all pending messages
    empty := false
    for !empty {
      select {
      case m := <-p.Transport.ControlRecv():
        messages = append(messages, m)
      default:
        empty = true
      }
    }

    // No new messages, so we can let somebody else run
    if len(messages) == 0 {
      runtime.Gosched()
      continue
    }

    // Process all pending messages
    any := false
    for _, m := range messages {
      if process(m) {
        any = true
      }
    }

    // No useful work done
    if !any {
      runtime.Gosched()
    }
  }
}

func (p *Peer) Bind(m interface {}) (any bool) {
  switch msg := m.(type) {
  case TrackerReq:
    any = true
    p.Transport.ControlSend(msg.from, TrackerRes{p.tracker})
  case Neighbours:
    any = true
    p.ids = msg.ids

    // Find if I'm a seed
    p.Transport.ControlSend(p.tracker, SeedReq{p.id})
  case SeedRes:
    any = true
    p.pieces = msg.pieces

    // Since we do Run here, it must be that it will not hang
    p.Run()
  default:
    if len(p.connectors) > 0 {
      any = true

      // All initialized
      p.RunRecv(m)
    } else {
      // Send message to myself
      p.Transport.ControlSend(p.id, m)
    }
  }
  return
}

func (p *Peer) Run() {
  // We want to bind all these variables here, so
  // we don't need any synchroization.
  fmt.Println(p.id, p.ids)

  // make per peer variables
  p.Storage   = NewStorage(p.id, p.pieces)
  p.Picker    = NewPicker(p.Storage)
  p.Transport = p.Transport
  p.Manager   = NewConnectionManager()
  p.Choker    = NewChoker(p.Manager, p.time)

  // make connectors
  for _, id := range p.ids {
    p.addConnector(id)
  }

  go p.Choker.Run()
}

func (p *Peer) RunRecv(m interface {}) {
  id   := ""

  // Redirect the message to the connector
  switch msg := m.(type) {
  case Choke:
    id = msg.id
    // fmt.Println("Msg:", p.id, reflect.TypeOf(msg), msg)
  case Unchoke:
    id = msg.id
    // fmt.Println("Msg:", p.id, reflect.TypeOf(msg), msg)
  case Interested:
    id = msg.id
    // fmt.Println("Msg:", p.id, reflect.TypeOf(msg), msg)
  case NotInterested:
    id = msg.id
    // fmt.Println("Msg:", p.id, reflect.TypeOf(msg), msg)
  case Have:
    id = msg.id
    // fmt.Println("Msg:", p.id, reflect.TypeOf(msg), msg)
  case Request:
    id = msg.id
    // fmt.Println("Msg:", p.id, reflect.TypeOf(msg), msg)
  case Piece:
    id = msg.id
    // fmt.Println("Msg:", p.id, reflect.TypeOf(msg), msg)
  case ConnReq:
    id = msg.id
    // fmt.Println("Msg:", p.id, reflect.TypeOf(msg), msg)
  }

  if id == "" {
    return
  }

  /**
   * The peer can have up to at most config.OutPeers connections
   * initiated by it.
   *
   * A perfect tracker would not keep "one-way" connections, that
   * is connections that are not reciprocated. (i.e. accepted as an
   * in-connection)
   *
   * To handle this, some protocols allow the peers to be
   * periodically changed. However, when using a perfect tracker, this
   * is not needed and we can set OutPeers = 0.
   */

  // if _, ok := p.connectors[id]; !ok && len(p.connectors) < maxPeers {
  if _, ok := p.connectors[id]; !ok {
    /*
     * This should not be reached when we having a perfect tracker.
     */
     p.addConnector(id)
  }

  if connector, ok := p.connectors[id]; ok {
    connector.Recv(m)
  }
}

func (p *Peer) addConnector(id string) {
  connector := NewConnector(p.id, id, p.Components)

  p.connectors[id] = connector
  p.Manager.AddConnector(connector)

  go connector.Run()
}
