package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

// A Forwarder is a structure used by a Leader to forward 'Have' messages
// towards all its Followers.
type Forwarder struct {
  *Leader

  from string
  to   string
}

func NewForwarder(leader *Leader, from, to string) *Forwarder {
  return &Forwarder{
    Leader : leader,

    from : from,
    to   : to,
  }
}

func (f *Forwarder) Recv(m interface {}) {
  id := f.GetId(m)

  // follower <-> leader <- peer

  if id == f.from {
    // follower -> leader

    switch msg := m.(type) {
    case torrent.Request:
      _, ok := f.Storage.Have(msg.Index)
      if !ok {
        // In this moment, the connection follower - leader is unchoked
        // but we don't have the piece. We will try to favor the transfer of
        // the piece by letting the picker know.

        picker := f.Picker.(*CachePicker)
        picker.GotRequest(msg.Index)
      }
    }
  } else {
    // leader <- peer

    switch msg := m.(type) {
    case torrent.Have:
      f.Transport.ControlSend(f.from, torrent.Have{
        Id    : f.Id,
        Index : msg.Index,
      })
    }
  }
}
