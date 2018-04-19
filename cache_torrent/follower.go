package cache_torrent

/**
 * A follower can:
 *   - upload to anybody
 *   - download only from same AS
 */

import (
  "github.com/danalex97/nfsTorrent/torrent"
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

func (l *Follower) Run() {
  fmt.Println("Follower running.")
}

func (l *Follower) Recv(m interface {}) {
  l.RunRecv(m, l.incomingConnection)
}

func (l *Follower) incomingConnection(id string) {
  if getAS(id) == getAS(l.Id) {
    // We make a bidirectional connection.
    torrent.
      NewConnector(l.Id, id, l.Components).
      WithHandshake().
      WithUpload().
      WithDownload().
      Register(l.Peer.Peer)
  } else {
    
  }
}
