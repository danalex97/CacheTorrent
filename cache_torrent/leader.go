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

  followerFwd map[string][]*Forwarder
  peerFwd     map[string][]*Forwarder
}

func NewLeader(p *Peer) *Leader {
  return &Leader{
    Peer : p,

    followerFwd : make(map[string][]*Forwarder),
    peerFwd     : make(map[string][]*Forwarder),
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
    if _, ok := l.Connectors[follower]; !ok {
      l.outgoingConnection(follower)
    }

    // Like this some connections from Leader to Leader may become
    // unidirectional [Race]
    if _, ok := l.Connectors[peer]; !ok {
      // If there is no connection with the Peer, we make a download
      // only connection. That is, we do no handshake and send a message
      // to the peer.

      torrent.
        NewConnector(l.Id, peer, l.Components).
        WithDownload().
        Register(l.Peer.Peer)
    }

    // Once the connections are made, we only need to register the forwarder
    fwd := NewForwarder(follower, peer)

    if _, ok := l.followerFwd[follower]; !ok {
       l.followerFwd[follower] = []*Forwarder{}
    }
    l.followerFwd[follower] = append(l.followerFwd[follower], fwd)

    if _, ok := l.peerFwd[peer]; !ok {
       l.peerFwd[peer] = []*Forwarder{}
    }
    l.peerFwd[peer] = append(l.peerFwd[peer], fwd)
  }

  l.Peer.RunRecv(m, l.incomingConnection)
}

func (l *Leader) outgoingConnection(id string) {
  // A leader has the initial protocol capabilities, that is
  // it's able to upload to anybody for outgoing connections.
  torrent.
    NewConnector(l.Id, id, l.Components).
    WithUpload().
    WithDownload().
    Register(l.Peer.Peer)
}

func (l *Leader) incomingConnection(id string) {
  l.outgoingConnection(id)
}
