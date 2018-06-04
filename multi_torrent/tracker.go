package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
  "github.com/danalex97/nfsTorrent/torrent"
)

type MultiTracker struct {
  election *MultiElection
  joined   map[string]bool

  *torrent.Tracker
}

func (t *MultiTracker) New(util TorrentNodeUtil) TorrentNode {
  tracker := new(MultiTracker)

  tracker.Tracker   = (tracker.Tracker.New(util)).(*torrent.Tracker)
  tracker.Transport = NewTransportProxy(tracker.Transport)

  tracker.election = NewMultiElection(
    MultiPeerMembers,
    tracker.Limit,
    tracker.Transport)
  tracker.joined   = make(map[string]bool)

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
    // We do not repeat the join messages
    id := ExternId(msg.Id)
    if _, ok := t.joined[id]; !ok {
      t.joined[id] = true
      t.Join(torrent.Join{id}, t.Neighbours)
    }
  default:
    t.Tracker.Recv(m)
  }

  t.election.Recv(m)
}

func (t *MultiTracker) Neighbours(id string) interface {} {
  ids := (t.Tracker.Neighbours(id)).(torrent.Neighbours).Ids

  t.election.NewJoin(id)

  // Keep message type compatible with cache_torrent protocol
  return cache_torrent.Neighbours{
    Ids : ids,
  }
}
