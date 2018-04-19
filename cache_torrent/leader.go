package cache_torrent

/**
 * A leader can:
 *  - download from anybody
 *  - upload to anybody(see race condition)
 * A leader can not:
 *  - upload to different AS via an indirect connection
 */

import (
  "github.com/danalex97/nfsTorrent/torrent"
  "fmt"
)

type Leader struct {
  *Peer
}

func NewLeader(p *Peer) *Leader {
  return &Leader{
    Peer : p,
  }
}

func (l *Leader) Run() {
  fmt.Println("Leader running.")
  l.Peer.Run(l.outgoingConnection)
}

func (l *Leader) Recv(m interface {}) {
  switch msg := m.(type) {
  case LeaderStart:
    follower := msg.Id
    peer     := msg.Dest

    // Make bidirectional connection to follower
    l.outgoingConnection(follower)

    // Like this some connections from Leader to Leader may become
    // unidirectional [Race]
    if _, ok := l.Connectors[peer]; !ok {
      // If there is no connection with the Peer, we make a download
      // only connection. That is, we do no handshake and send a message
      // to the peer.


    }


  }

  l.Peer.RunRecv(m, l.incomingConnection)
}

func (l *Leader) outgoingConnection(id string) {
  // A leader has the initial protocol capabilities, that is
  // it's able to upload to anybody for outgoing connections.
  torrent.
    NewConnector(l.Id, id, l.Components).
    WithHandshake().
    WithUpload().
    WithDownload().
    Register(l.Peer.Peer)
}

func (l *Leader) incomingConnection(id string) {
  l.outgoingConnection(id)
}
