package bias_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "github.com/danalex97/nfsTorrent/config"

  "strings"
  "math/rand"
  "math"
)

var KPercent config.Const = config.NewConst(config.Bias)

type Tracker struct {
  asId map[string][]string

  *torrent.Tracker
}

func (t *Tracker) New(util TorrentNodeUtil) TorrentNode {
  tracker := new(Tracker)
  tracker.asId     = make(map[string][]string)
  tracker.Tracker  = (tracker.Tracker.New(util)).(*torrent.Tracker)
  return tracker
}

func (t *Tracker) OnJoin() {
  go t.CheckMessages(t.Recv)
}

func (t *Tracker) Recv(m interface {}) {
  switch msg := m.(type) {
  /* New Protocol. */
  case torrent.Join:
    id := msg.Id
    as := getAS(id)
    if _, ok := t.asId[as]; !ok {
      t.asId[as] = []string{}
    }
    t.asId[as] = append(t.asId[as], id)

    t.Join(msg, t.Neighbours)
  default:
    /* Backward compatibility. */
    t.Tracker.Recv(m)
  }
}

func getAS(id string) string {
  // We assume that ID is of form [AS].[NBR]
  return strings.Split(id, ".")[0]
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
  neighbours := torrent.Neighbours{[]string{}}
  as := getAS(id)
  B  := t.Neigh - int(math.Floor(float64(KPercent.Int()) * float64(t.Neigh) / 100))

  if nbr, _ :=t.asId[as]; B > len(nbr) {
    B = len(nbr) - 1
  }

  ids := []string{}
  for _, newId := range t.asId[as] {
    if newId != id {
      ids = append(ids, newId)
    }
  }

  for i := 0; i < t.Neigh; i++ {
    if i >= B {
      ids = t.Ids
    }
    id := newNeigh(ids, append(neighbours.Ids, id))
    neighbours.Ids = append(neighbours.Ids, id)
  }
  return neighbours
}
