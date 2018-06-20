package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
  "math/rand"
)

// A Follower is a Peer with limited capabilities.
// A Follower can:
//   - upload to anybody
//   - download only from same AS
//
// The Follower mostly acts like a BitTorrent peer, whereas the Leader would
// have further modifications such as message forwarding and special
// piece picking, basically acting as a proxy between followers and other peers.
//
// The same autonomous system connections are establised as in the BitTorrent
// protocol. For indirect connections, the Follower sends the 'leaderStart'
// message towards randomly-chosen leaders of its autonomous system. When the
// Leader receives the leaderStart message, it opens a download-only connection
// with the remote peer that the follower wants to communicate with.
type Follower struct {
  *Peer
}

func NewFollower(p *Peer) *Follower {
  return &Follower{
    Peer : p,
  }
}

func (f *Follower) Run() {
  f.Peer.Run(f.outgoingConnection)
}

func (f *Follower) Recv(m interface {}) {
  f.Peer.RunRecv(f.GetId(m), m, f.incomingConnection)
}

func (f *Follower) outgoingConnection(id string) {
  // Outgoing connections are handled similarly to incoming
  // connections.
  f.incomingConnection(id)
}

func (f *Follower) incomingConnection(id string) {
  if getAS(id) == getAS(f.Id) {
    // We make a bidirectional connection.
    torrent.
      NewConnector(f.Id, id, f.Components).
      WithUpload(NewUpload).
      WithDownload(NewDownload).
      Register(f.Peer.Peer)
  } else {
    // We receive an incoming connection from a diffrent AS

    // 1. Open an upload connection towards that node
    torrent.
      NewConnector(f.Id, id, f.Components).
      WithUpload(NewUpload).
      Register(f.Peer.Peer)

    // 2. Open an indirect connection
    leader := f.Leaders[rand.Intn(len(f.Leaders))]
    f.openIndirect(leader, id)
  }
}

func (f *Follower) openIndirect(leader, target string) {
  if _, ok := f.Connectors[leader]; !ok {
    // We have no direct connection with the leader, so we make one
    torrent.
      NewConnector(f.Id, leader, f.Components).
      WithUpload(NewUpload).
      WithDownload(NewDownload).
      Register(f.Peer.Peer)
  }

  // Now we let the Leader know we want an indirect connection.
  f.Transport.ControlSend(leader, LeaderStart{
    Id   : f.Id,
    Dest : target,
  })
}
