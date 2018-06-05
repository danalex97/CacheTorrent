package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "github.com/danalex97/nfsTorrent/log"

  "runtime"
  "reflect"
  // "fmt"
)

var inPeers  config.Const = config.NewConst(config.InPeers)
var outPeers config.Const = config.NewConst(config.OutPeers)
var progress config.Const = config.NewConst(config.AllNodesRun)

type Peer struct {
  *Components

  // Information used for initializing the peer.
  Id      string
  Tracker string
  Ids     []string

  // -- BitTorrent protocol --
  Pieces []PieceMeta

  // -- BitTorrent components --
  // A map of Connectiors for all the nodes provided by the Tracker or by other
  // peers via connection initialization.
  Connectors  map[string]Runner

  // Used only to identify tracker via connecting a bootstrap node which
  // follows the BitTorrent protocol.
  join string

  // Used to keep if logging is on. Logs can be visualized via the vizualtion
  // tool in 'visual' folder.
  logging bool

  // Used to support multiple message arrival(which indirectly offers support
  // for MultiTorrents)
  bound bool
}

type Binder    func(m interface {}) int
type Processor func(m interface {}, state int)
type ConnAdder func(id string)

// The Bind states of a peer are:
//  - BindNone - the initial state during the bind stage; if the message format
//     is unexpected the state will remain BindNode
//  - BindDone - the state which announced that work was done, even though the
//     work is done only for building the peer's connections
//  - BindRecv - if the initialization phase is done(i.e. the peer began to
//     build its connection list) and the message can be passed towards a
//     specific Connection based on the sender's ID
//  - BindRun - the BindRun state anounced that the phase Tracker sent all the
//     necessary informations so we can finish the initialization phase; when
//     the peer arrives in the BindRun state the function Run is called
//     synchronously
const (
  BindNone = iota
  BindDone = iota
  BindRecv = iota
  BindRun  = iota
)

// The OnJoin function is called when the peer join the system and is part of
// the Implementation of Torrent Node interface.
func (p *Peer) OnJoin() {
  // If the Transport is missing, it must be we are
  // on torrent-less node
  if p.Transport == nil {
    return
  }

  // We need to make this unblocking for latency support.
  go func() {
    p.Init()
    go p.CheckMessages(p.Bind, p.Process)
  }()
}

// The OnLeave function is called when the peer leaves0 the system and is
// part of the Implementation of Torrent Node interface.
func (p *Peer) OnLeave() {
}

// The New function is the construction function for the TorrentNode. The
// peer should be partially initialized in this function by using the
// TorrentNodeUtil's components. The TorrentNodeUtil is an interface which
// allows a node to interact with the Speer simulator.
func (p *Peer) New(util TorrentNodeUtil) TorrentNode {
  peer := new(Peer)

  peer.Id         = util.Id()
  peer.join       = util.Join()
  peer.Components = new(Components)
  peer.Transport  = util.Transport()

  peer.Pieces     = []PieceMeta{}
  peer.Connectors = make(map[string]Runner)

  peer.Time = util.Time()

  peer.bound = false

  return peer
}

// Initilization function called synchronously after a node has joined the
// peer-to-peer system. In this function we find who the tracker is by
// interacting with the bootstrap node, respond to other peers which ask
// who the tracker is and, finally, send the tracker a Join request.
func (p *Peer) Init() {
  // Find out who the tracker is.
  p.Transport.ControlSend(p.join, TrackerReq{p.Id})

  queue    := []TrackerReq{}
  response := false
  for !response {
    m := <-p.Transport.ControlRecv()
    switch msg := m.(type) {
    case TrackerRes:
      p.Tracker = msg.Id
      response  = true
    case TrackerReq:
      // We may receive TrackerReq while waiting for our own response.
      queue = append(queue, msg)
    }
  }

  // Send all the responses for TrackerReq to other peers.
  for _, msg := range queue {
    p.Transport.ControlSend(msg.From, TrackerRes{p.Tracker})
  }

  // The peer should be initialized
  log.Println("Node", p.Id, "started with tracker", p.Tracker)

  // Send join message to the tracker
  p.Transport.ControlSend(p.Tracker, Join{p.Id})

  // Initialize deatailed logging.
  p.logging = log.HasLogfile();
}

// CheckMessages is a function which runs until the node departs the system.
// It registers all the incoming messages. The messages are checked using a
// Binder which identifies the state of the Peer and the messages are processed
// by being passed to a Processor. The Processor looks at the state returned by
// the binder and acts correspondingly dispaching the message accordingly.
//
// The CheckMessages function is exposed in the public API to make building on
// the current protocol easier by following the same Binder-Processor pattern.
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
      // fmt.Println(reflect.TypeOf(m).String(), m)

      state := bind(m)
      if state != BindNone {
        any = true
      }

      if state == BindRun {
        if !p.bound {
          // Notify the progress properties.
          progress.Ref().(*config.WGProgress).Add()
          p.bound = true
        }
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

// Process is a Processor function which calls the Run function when the node
// reaches the BindRun state. Afterwards, when the node is in the BindRecv
// state, the Process function will redirect the message to the RunRecv method.
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

// The run function is called at the end of the initialization phase. It
// initializes most of the internal components of a peer exposed by the
// Components structure. Moreover, it starts all the asynchronous routines such
// as the Choker and registers all the connections provided by the Tracker.
func (p *Peer) Run(connAdd ConnAdder) {
  // We want to bind all these variables here, so
  // we don't need any synchroization.
  // log.Println(p.Id, p.Ids)

  // Make per peer components.
  p.Storage   = NewStorage(p.Id, p.Pieces, p.Time)
  p.Picker    = NewPicker(p.Storage)
  // p.Transport = p.Transport
  p.Manager   = NewConnectionManager()
  p.Choker    = NewChoker(p.Manager, p.Time)

  // Make connectors.
  for _, id := range p.Ids {
    // Defensive programming.
    if _, ok := p.Connectors[id]; !ok {
      connAdd(id)
    }
  }

  go p.Choker.Run()
}

// A function used to extract the Id from an incoming message.
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

// The RunRecv function is used to dispach the message towards a
// Connector's Recv method.
func (p *Peer) RunRecv(id string, m interface {}, connAdd ConnAdder) {
  if id == "" {
    return
  }

  // The peer can have up to at most config.Config.OutPeers connections
  // initiated by it.
  //
  // A perfect tracker would not keep "one-way" connections, that
  // is connections that are not reciprocated. (i.e. accepted as an
  // in-connection)
  //
  // To handle this, some protocols allow the peers to be
  // periodically changed. However, when using a perfect tracker, this
  // is not needed and we can set OutPeers = 0.

  // if _, ok := p.Connectors[id]; !ok && len(p.Connectors) < inPeers + outPeers {
  if _, ok := p.Connectors[id]; !ok {
    // This should not be reached when we having a perfect tracker.
    connAdd(id)
  }

  if connector, ok := p.Connectors[id]; ok {
    connector.Recv(m)
  }
}

// Function called when a new Connector towards the node id should be
// initialized.
func (p *Peer) AddConnector(id string) {
  NewConnector(p.Id, id, p.Components).
  WithUpload(NewUpload).
  WithDownload(NewDownload).
  Register(p)
}
