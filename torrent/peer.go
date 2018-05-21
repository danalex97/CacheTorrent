package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "github.com/danalex97/nfsTorrent/log"

  "runtime"
  "reflect"
)

var inPeers  config.Const = config.NewConst(config.InPeers)
var outPeers config.Const = config.NewConst(config.OutPeers)
var progress config.Const = config.NewConst(config.AllNodesRun)

type Peer struct {
  *Components

  Id      string
  Tracker string
  Ids     []string

  // BitTorrent protocol
  Pieces []PieceMeta

  // BitTorrent components
  Connectors  map[string]Runner // the connectors that were chosen by tracker

  // used only to identify tracker
  join    string

  // used to keep if logging is on
  logging bool
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
  log.Println("Node", p.Id, "started with tracker", p.Tracker)

  // Send join message to the tracker
  p.Transport.ControlSend(p.Tracker, Join{p.Id})

  // Initialize deatailed logging.
  p.logging = log.HasLogfile();
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
      // Notify progress properties.
      progress.Ref().(*config.WGProgress).Progress(p.Id)

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

      if state == BindRun {
        // Notify the progress properties.
        progress.Ref().(*config.WGProgress).Add()
      }

      process(m, state)

      // Log messages.
      if p.logging {
        log.LogPacket(log.Packet{
          Src  : reflect.ValueOf(m).FieldByName("Id").String(),
          Dst  : p.Id,
          Type : reflect.TypeOf(m).String(),
        })
      }
    }

    // Notify progress properties.
    progress.Ref().(*config.WGProgress).Progress(p.Id)

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
    p.RunRecv(p.GetId(m), m, p.AddConnector)
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
    if len(p.Pieces) > 0 {
      // I am seed
      log.Println("Seed upload:", p.Transport.Up())
    }
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
  // log.Println(p.Id, p.Ids)

  // make per peer variables
  p.Storage   = NewStorage(p.Id, p.Pieces, p.Time)
  p.Picker    = NewPicker(p.Storage)
  p.Transport = p.Transport
  p.Manager   = NewConnectionManager()
  p.Choker    = NewChoker(p.Manager, p.Time)

  // Make connectors
  for _, id := range p.Ids {
    // Defensive programming
    if _, ok := p.Connectors[id]; !ok {
      connAdd(id)
    }
  }

  go p.Choker.Run()
}

func (p *Peer) GetId(m interface {}) (id string) {
  switch msg := m.(type) {
  case Choke:
    id = msg.Id
  case Unchoke:
    id = msg.Id
  case Interested:
    id = msg.Id
  case NotInterested:
    id = msg.Id
  case Have:
    id = msg.Id
  case Request:
    id = msg.Id
  case Piece:
    id = msg.Id
  case ConnReq:
    id = msg.Id
  }
  return
}

func (p *Peer) RunRecv(id string, m interface {}, connAdd ConnAdder) {
  if id == "" {
    return
  }

  /**
   * The peer can have up to at most config.Config.OutPeers connections
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

  // if _, ok := p.Connectors[id]; !ok && len(p.Connectors) < inPeers + outPeers {
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
  WithUpload(NewUpload).
  WithDownload(NewDownload).
  Register(p)
}
