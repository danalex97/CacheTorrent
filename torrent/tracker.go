package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "runtime"
  "math/rand"
)

const MinNodes       int = config.MinNodes
const PeerNeighbours int = config.OutPeers
const Seeds          int = config.Seeds
const Pieces         int = config.Pieces
const PieceSize      int = config.PieceSize

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

  tracker.Limit = MinNodes
  tracker.Neigh = PeerNeighbours
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
      if i < Seeds {
        // It's a seed
        ps     := []PieceMeta{}
        begin  := 0
        length := PieceSize

        for j := 0; j < Pieces; j++ {
          ps    = append(ps, PieceMeta{j, begin, length})
          begin = begin + length
        }

        return SeedRes{ps}
      } else {
        return SeedRes{[]PieceMeta{}}
      }
    }
  }
  panic("Seed request: from ID not found!")
}
