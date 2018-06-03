package bias_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
  "github.com/danalex97/nfsTorrent/torrent"
)

// The CacheTracker is component which allows usage of bias neighbour election
// together with the CacheTorrent protocol.
type CacheTracker struct {
  election *cache_torrent.Election

  *Tracker
}

func (t *CacheTracker) New(util TorrentNodeUtil) TorrentNode {
  tracker := new(CacheTracker)
  tracker.Tracker  = (tracker.Tracker.New(util)).(*Tracker)
  tracker.election = cache_torrent.NewElection(tracker.Limit, tracker.Transport)
  return tracker
}

func (t *CacheTracker) OnJoin() {
  go t.election.Run()
  go t.CheckMessages(t.Recv)
}

func (t *CacheTracker) Recv(m interface {}) {
  switch msg := m.(type) {
  case torrent.Join:
    // Override the Neighbours method so we can send the notify the election.
    t.Join(msg, t.Neighbours)
  default:
    t.Tracker.Recv(m)
  }

  t.election.Recv(m)
}

func (t *CacheTracker) Neighbours(id string) interface {} {
  ids := (t.Tracker.Neighbours(id)).(torrent.Neighbours).Ids

  // Notify the election.
  t.election.NewJoin(id)

  // Keep message type compatible with cache_torrent protocol
  return cache_torrent.Neighbours{
    Ids : ids,
  }
}
