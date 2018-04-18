package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
)

type Tracker struct {
  election *Election

  *torrent.Tracker
}

func (t *Tracker) New(util TorrentNodeUtil) TorrentNode {
  tracker := new(Tracker)
  tracker.Tracker  = (tracker.Tracker.New(util)).(*torrent.Tracker)
  tracker.election = NewElection(tracker.Limit, tracker.Transport)
  return tracker
}

func (t *Tracker) OnJoin() {
  go t.election.Run()
  go t.CheckMessages(t.Recv)
}

func (t *Tracker) Recv(m interface {}) {
  switch msg := m.(type) {
  /* New Protocol. */
  case torrent.Join:
    t.Join(msg, t.Neighbours)
  default:
    /* Backward compatibility. */
    t.Tracker.Recv(m)
  }

  /* New Protocol. */
  t.election.Recv(m)
}

func (t *Tracker) Neighbours(id string) interface {} {
  ids := (t.Tracker.Neighbours(id)).(torrent.Neighbours).Ids

  t.election.NewJoin(id)

  return Neighbours{
    Ids : ids,
  }
}

func (t *Tracker) Local(id string) []string {
  local := []string{}
  for _, nid := range t.Ids {
    if getAS(id) == getAS(nid) && id != nid {
      local = append(local, nid)
    }
  }
  return local
}
