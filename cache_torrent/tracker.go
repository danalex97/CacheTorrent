package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "strings"
  "sync"
)

type Tracker struct {
  sync.Mutex
  leaders    map[string][]string

  *torrent.Tracker
}

func (t *Tracker) New(util TorrentNodeUtil) TorrentNode {
  tracker := new(Tracker)
  tracker.Tracker = (tracker.Tracker.New(util)).(*torrent.Tracker)
  return tracker
}

func (t *Tracker) OnJoin() {
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
}

func (t *Tracker) Neighbours(id string) interface {} {
  ids := (t.Tracker.Neighbours(id)).(torrent.Neighbours).Ids

  // Run Leader election
  go t.Elect(id, t.Local(id))

  return Neighbours{
    Ids : ids,
  }
}

func getAS(id string) string {
  // We assume that ID is of form [AS].[NBR]
  return strings.Split(id, ".")[0]
}

func (t *Tracker) Elect(id string, ids []string) {
  as := getAS(id)

  // Run leader election for an AS
  t.Lock()
  defer t.Unlock()

  // We don't handle ulterior joins for the moment
  if leaders, ok := t.leaders[as]; ok {
    // If election for AS already took place leet peer know
    t.Transport.ControlSend(id, Leaders{leaders})
    return
  }

  // We need to do the election

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
