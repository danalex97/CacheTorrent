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
  switch m.(type) {
  default:
    /* Backward compatibility. */
    t.Tracker.Recv(m)
  }
}
