package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "runtime"
  "fmt"
)

const MaxPeers int = config.InPeers + config.OutPeers

type Peer struct {
  *Components

  Id      string
  Tracker string
  Ids     []string

  Time      func() int

  // BitTorrent protocol
  Pieces []PieceMeta

  // BitTorrent components
  Connectors  map[string]Runner // the connectors that were chosen by tracker

  // used only to identify tracker
  join    string
}

type Binder    func(m interface {}) int
type Processor func(m interface {}, state int)
type ConnAdder func(id string)

const (
  BindNone = iota
  BindDone = iota
  BindRecv = iota
  BindRun  = iota
)

/* Implementation of Torrent Node interface. */
func (p *Peer) OnJoin() {
  // If the Transport is missing, it must be we are
  // on torrent-less node
  if p.Transport == nil {
    return
  }

  p.Init()
  go p.CheckMessages(p.Bind, p.Process)
}

func (p *Peer) OnLeave() {
}

func (p *Peer) New(util TorrentNodeUtil) TorrentNode {
  peer := new(Peer)

  peer.Id        = util.Id()
  peer.join      = util.Join()
  peer.Components = new(Components)
  peer.Transport  = util.Transport()

  peer.Pieces     = []PieceMeta{}
  peer.Connectors = make(map[string]Runner)

  peer.Time = util.Time()

  return peer
}

/* Internal functions. */
func (p *Peer) Init() {
  // Find out who the tracker is
  p.Transport.ControlSend(p.join, TrackerReq{p.Id})

  msg := <-p.Transport.ControlRecv()
  p.Tracker = msg.(TrackerRes).Id

  // The peer should be initialized
  fmt.Printf("Node %s started with tracker %s\n", p.Id, p.Tracker)

  // Send join message to the tracker
  p.Transport.ControlSend(p.Tracker, Join{p.Id})
}

func (p *Peer) CheckMessages(bind Binder, process Processor) {
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
      state := bind(m)
      if state != BindNone {
        any = true
      }
      process(m, state)
    }

    // No useful work done
    if !any {
      runtime.Gosched()
    }
  }
}

func (p *Peer) Process(m interface {}, state int) {
  switch state {
  case BindRun:
    p.Run(p.AddConnector)
  case BindRecv:
    p.RunRecv(m, p.AddConnector)
  }
}

func (p *Peer) Bind(m interface {}) (state int) {
  state = BindNone

  switch msg := m.(type) {
  case TrackerReq:
    state = BindDone

    p.Transport.ControlSend(msg.From, TrackerRes{p.Tracker})
  case Neighbours:
    state = BindDone
    p.Ids = msg.Ids

    // Find if I'm a seed
    p.Transport.ControlSend(p.Tracker, SeedReq{p.Id})
  case SeedRes:
    state = BindRun

    p.Pieces = msg.Pieces
  default:
    if len(p.Connectors) > 0 {
      state = BindRecv
    } else {
      state = BindNone

      // Send message to myself
      p.Transport.ControlSend(p.Id, m)
    }
  }

  return
}

func (p *Peer) Run(connAdd ConnAdder) {
  // We want to bind all these variables here, so
  // we don't need any synchroization.
  fmt.Println(p.Id, p.Ids)

  // make per peer variables
  p.Storage   = NewStorage(p.Id, p.Pieces)
  p.Picker    = NewPicker(p.Storage)
  p.Transport = p.Transport
  p.Manager   = NewConnectionManager()
  p.Choker    = NewChoker(p.Manager, p.Time)

  // make connectors
  for _, id := range p.Ids {
    connAdd(id)
  }

  go p.Choker.Run()
}

func (p *Peer) GetId(m interface {}) (id string){
  switch msg := m.(type) {
  case Choke:
    id = msg.Id
    // fmt.Println("Msg:", p.Id, reflect.TypeOf(msg), msg)
  case Unchoke:
    id = msg.Id
    // fmt.Println("Msg:", p.Id, reflect.TypeOf(msg), msg)
  case Interested:
    id = msg.Id
    // fmt.Println("Msg:", p.Id, reflect.TypeOf(msg), msg)
  case NotInterested:
    id = msg.Id
    // fmt.Println("Msg:", p.Id, reflect.TypeOf(msg), msg)
  case Have:
    id = msg.Id
    // fmt.Println("Msg:", p.Id, reflect.TypeOf(msg), msg)
  case Request:
    id = msg.Id
    // fmt.Println("Msg:", p.Id, reflect.TypeOf(msg), msg)
  case Piece:
    id = msg.Id
    // fmt.Println("Msg:", p.Id, reflect.TypeOf(msg), msg)
  case ConnReq:
    id = msg.Id
    // fmt.Println("Msg:", p.Id, reflect.TypeOf(msg), msg)
  }
  return
}

func (p *Peer) RunRecv(m interface {}, connAdd ConnAdder) {
  id := p.GetId(m)

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

  // if _, ok := p.Connectors[id]; !ok && len(p.Connectors) < maxPeers {
  if _, ok := p.Connectors[id]; !ok {
    /*
     * This should not be reached when we having a perfect tracker.
     */
     connAdd(id)
  }

  if connector, ok := p.Connectors[id]; ok {
    connector.Recv(m)
  }
}

func (p *Peer) AddConnector(id string) {
  NewConnector(p.Id, id, p.Components).
  WithUpload().
  WithDownload().
  Register(p)
}
