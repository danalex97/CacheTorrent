package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "runtime"
  "math/rand"
)

var MinNodes       config.Const = config.NewConst(config.MinNodes)
var PeerNeighbours config.Const = config.NewConst(config.OutPeers)
var Seeds          config.Const = config.NewConst(config.Seeds)
var Pieces         config.Const = config.NewConst(config.Pieces)
var PieceSize      config.Const = config.NewConst(config.PieceSize)

// The Tracker works as a bootstrapping mechanism. Firstly, when a peer wants
// to join the network it will ask the tracker for a list of peers to which it
// can connect to. The tracker has a pool of peers, from which it randomly
// chooses a fixed number of peers to pass to the requester. The second step,
// thus, consists of the tracker sending the list of peers to the requester.
// The peer will then try to connect to any of these peers.
//
// A modification that we make from the original Tracker is the fact that we
// use SeedReq for intergration with Speer. For simplicity, a Tracker will
// decide which nodes are seeds. This messages are useless in a real deployment.
//
// The Tracker reacts to messages:
//  - TrackerReq
//  - Join
//  - SeedReq
type Tracker struct {
  Ids       []string

  Limit int
  Neigh int
  Id    string

  Transport Transport
}

func (t *Tracker) New(util TorrentNodeUtil) TorrentNode {
  tracker := new(Tracker)

  tracker.Ids   = []string{}

  tracker.Limit = MinNodes.Int()
  tracker.Neigh = PeerNeighbours.Int()
  tracker.Id    = util.Id()

  tracker.Transport = util.Transport()

  return tracker
}

func (t *Tracker) OnJoin() {
  go t.CheckMessages(t.Recv)
}

func (t *Tracker) OnLeave() {
}

func (t *Tracker) CheckMessages(process func(interface {})) {
  for {
    select {
    case m, ok := <-t.Transport.ControlRecv():
      if !ok {
        // Channel closed
        break
      }

      process(m)
    default:
      // allow other nodes in simulation run
      runtime.Gosched()
    }
  }
}

func (t *Tracker) Recv(m interface {}) {
  switch msg := m.(type) {
  case Join:
    t.Join(msg, t.Neighbours)
  case TrackerReq:
    t.Transport.ControlSend(msg.From, TrackerRes{t.Id})
  case SeedReq:
    t.Transport.ControlSend(msg.From, t.seedRequest(msg))
  }
}

func (t *Tracker) Join(msg Join, getNeighbours func(string) interface {}) {
  // Add the new peer
  t.Ids = append(t.Ids, msg.Id)

  // We need to let all pending nodes know the
  // neighbour list
  if len(t.Ids) == t.Limit {
    for _, id := range t.Ids {
      t.Transport.ControlSend(id, getNeighbours(id))
    }
  }

  // We need to let the current node know the
  // neighbour list
  if len(t.Ids) > t.Limit {
    id := msg.Id
    t.Transport.ControlSend(id, getNeighbours(id))
  }
}

func newNeigh(allIds, ids []string) string {
  n := allIds[rand.Intn(len(allIds))]
  for _, i := range ids {
    if i == n {
      return newNeigh(allIds, ids)
    }
  }
  return n
}

func (t *Tracker) Neighbours(id string) interface {} {
  neighbours := Neighbours{[]string{}}
  for i := 0; i < t.Neigh; i++ {
    id := newNeigh(t.Ids, append(neighbours.Ids, id))
    neighbours.Ids = append(neighbours.Ids, id)
  }
  return neighbours
}

func (t *Tracker) seedRequest(req SeedReq) SeedRes {
  for i, id := range t.Ids {
    if id == req.From {
      if i < Seeds.Int() {
        // It's a seed
        ps     := []PieceMeta{}
        begin  := 0
        length := PieceSize.Int()

        for j := 0; j < Pieces.Int(); j++ {
          ps    = append(ps, PieceMeta{j, begin, length})
          begin = begin + length
        }

        return SeedRes{ps}
      } else {
        return SeedRes{[]PieceMeta{}}
      }
    }
  }
  panic("Seed request: from " + req.From + " not found!")
}
