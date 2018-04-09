package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "runtime"
  "math/rand"
)

const minNodes       int = config.MinNodes
const peerNeighbours int = config.PeerNeighbours
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
  tracker.id        = util.Id()

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
        case join:
          t.join(msg)
        case trackerReq:
          t.transport.ControlSend(msg.from, trackerRes{t.id})
        case seedReq:
          t.transport.ControlSend(msg.from, t.seedRequest(msg))
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

func (t *Tracker) join(msg join) {
  // Add the new peer
  t.ids = append(t.ids, msg.id)

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
    id := msg.id
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

func (t *Tracker) neighbours(id string) neighbours {
  neighbours := neighbours{[]string{}}
  for i := 0; i < t.neigh; i++ {
    id := newNeigh(t.ids, append(neighbours.ids, id))
    neighbours.ids = append(neighbours.ids, id)
  }
  return neighbours
}

func (t *Tracker) seedRequest(req seedReq) seedRes {
  for i, id := range t.ids {
    if id == req.from {
      if i < seeds {
        // It's a seed
        ps     := []request{}
        begin  := 0
        length := pieceSize

        for j := 0; j < pieces; j++ {
          ps    = append(ps, request{j, begin, length})
          begin = begin + length
        }

        return seedRes{ps}
      } else {
        return seedRes{[]request{}}
      }
    }
  }
  panic("Seed request: from ID not found!")
}
