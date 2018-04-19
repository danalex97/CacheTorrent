package cache_torrent

/**
 * A follower can:
 *   - upload to anybody
 *   - download only from same AS
 */

import (
  "github.com/danalex97/nfsTorrent/torrent"
  "math/rand"
  "fmt"
)

type Follower struct {
  *Peer
}

func NewFollower(p *Peer) *Follower {
  return &Follower{
    Peer : p,
  }
}

func (f *Follower) Run() {
  fmt.Println("Follower running.")
  f.Peer.Run(f.outgoingConnection)
}

func (f *Follower) Recv(m interface {}) {
  f.Peer.RunRecv(m, f.incomingConnection)
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
      WithHandshake().
      WithUpload().
      WithDownload().
      Register(f.Peer.Peer)
  } else {
    // We receive an incoming connection from a diffrent AS

    // 1. Open an upload connection towards that node
    torrent.
      NewConnector(f.Id, id, f.Components).
      WithHandshake().
      WithUpload().
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
      WithHandshake().
      WithUpload().
      WithDownload().
      Register(f.Peer.Peer)
  }

  // Now we let the Leader know we want an indirect connection.
  f.Transport.ControlSend(leader, LeaderStart{
    Id   : f.Id,
    Dest : target,
  })
}
