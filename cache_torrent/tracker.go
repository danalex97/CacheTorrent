package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
)

type Tracker struct {
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
  return Neighbours{
    Ids   : (t.Tracker.Neighbours(id)).(torrent.Neighbours).Ids,
    Local : t.Local(id),
  }
}

func (t *Tracker) Local(id string) []string {
  return []string{}
}
