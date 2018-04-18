package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "runtime"
  "math/rand"
)

const minNodes       int = config.MinNodes
const peerNeighbours int = config.OutPeers
const seeds          int = config.Seeds
const pieces         int = config.Pieces
const pieceSize      int = config.PieceSize

type Tracker struct {
  ids       []string

  limit int
  neigh int
  id    string

  transport Transport
}

func (t *Tracker) New(util TorrentNodeUtil) TorrentNode {
  tracker := new(Tracker)

  tracker.ids       = []string{}

  tracker.limit = minNodes
  tracker.neigh = peerNeighbours
  tracker.id     = util.Id()

  tracker.transport = util.Transport()

  return tracker
}

func (t *Tracker) OnJoin() {
  go func() {
    for {
      select {
      case m, ok := <-t.transport.ControlRecv():
        if !ok {
          continue
        }

        switch msg := m.(type) {
        case Join:
          t.join(msg)
        case TrackerReq:
          t.transport.ControlSend(msg.From, TrackerRes{t.id})
        case SeedReq:
          t.transport.ControlSend(msg.From, t.seedRequest(msg))
        }

      default:
        // allow other nodes in simulation run
        runtime.Gosched()
      }
    }
  }()
}

func (t *Tracker) OnLeave() {
}

func (t *Tracker) join(msg Join) {
  // Add the new peer
  t.ids = append(t.ids, msg.Id)

  // We need to let all pending nodes know the
  // neighbour list
  if len(t.ids) == t.limit {
    for _, id := range t.ids {
      t.transport.ControlSend(id, t.neighbours(id))
    }
  }

  // We need to let the current node know the
  // neighbour list
  if len(t.ids) > t.limit {
    id := msg.Id
    t.transport.ControlSend(id, t.neighbours(id))
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

func (t *Tracker) neighbours(id string) Neighbours {
  neighbours := Neighbours{[]string{}}
  for i := 0; i < t.neigh; i++ {
    id := newNeigh(t.ids, append(neighbours.Ids, id))
    neighbours.Ids = append(neighbours.Ids, id)
  }
  return neighbours
}

func (t *Tracker) seedRequest(req SeedReq) SeedRes {
  for i, id := range t.ids {
    if id == req.From {
      if i < seeds {
        // It's a seed
        ps     := []PieceMeta{}
        begin  := 0
        length := pieceSize

        for j := 0; j < pieces; j++ {
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
