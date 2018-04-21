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

func (l *Leader) GetId(m interface {}) string {
  switch msg := m.(type) {
  case LeaderStart:
    return msg.Id
  default:
    return l.Peer.GetId(m)
  }
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

    if _, ok := l.Connectors[peer]; !ok {
      // If there is no connection with the Peer, we make a download
      // only connection. That is, we do no handshake and send a message
      // to the peer.

      // We can add the upload component since the other peer does
      // upload only if it's a follower, so our Upload will do nothing
      // since it will be always choked.
      torrent.
        NewConnector(l.Id, peer, l.Components).
        WithDownload().
        WithUpload().
        Register(l.Peer.Peer)
    }

    // Once the connections are made, we only need to register the forwarder
    l.registerForwarder(follower, peer)
  }

  // Send the messages to corresponding forwarders
  l.forward(m)

  l.Peer.RunRecv(m, l.incomingConnection)
}

func (l *Leader) forward(m interface {}) {
  id := l.GetId(m)

  forward := func(mp map[string][]*Forwarder) {
    if _, ok := mp[id]; ok {
      for _, fwd := range mp[id] {
        fwd.Recv(m)
      }
    }
  }

  forward(l.followerFwd)
  forward(l.peerFwd)
}

func (l *Leader) registerForwarder(follower, peer string) {
  fwd := NewForwarder(l, follower, peer)

  if _, ok := l.followerFwd[follower]; !ok {
     l.followerFwd[follower] = []*Forwarder{}
  }
  l.followerFwd[follower] = append(l.followerFwd[follower], fwd)

  if _, ok := l.peerFwd[peer]; !ok {
     l.peerFwd[peer] = []*Forwarder{}
  }
  l.peerFwd[peer] = append(l.peerFwd[peer], fwd)
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
