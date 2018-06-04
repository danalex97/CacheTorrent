package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
  "github.com/danalex97/nfsTorrent/torrent"
)

type MultiTracker struct {
  election *cache_torrent.Election

  *torrent.Tracker
}

func (t *MultiTracker) New(util TorrentNodeUtil) TorrentNode {
  tracker := new(MultiTracker)

  tracker.Tracker   = (tracker.Tracker.New(util)).(*torrent.Tracker)
  tracker.Transport = NewTransportProxy(tracker.Transport)

  tracker.election = cache_torrent.NewElection(tracker.Limit, tracker.Transport)

  return tracker
}

func (t *MultiTracker) OnJoin() {
  go t.election.Run()
  go t.CheckMessages(t.Recv)
}

func (t *MultiTracker) Recv(m interface {}) {
  switch msg := m.(type) {
  case torrent.Join:
    // We ignore torrent.Join messages
  case Join:
    // Override the Neighbours method so we can send the notify the election.
    t.Join(torrent.Join{msg.Id}, t.Neighbours)
  default:
    t.Tracker.Recv(m)
  }

  t.election.Recv(m)
}

func (t *MultiTracker) Neighbours(id string) interface {} {
  ids := (t.Tracker.Neighbours(id)).(torrent.Neighbours).Ids

  // Notify the election.
  t.election.NewJoin(id)

  // Keep message type compatible with cache_torrent protocol
  return cache_torrent.Neighbours{
    Ids : ids,
  }
}
