package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
  "github.com/danalex97/nfsTorrent/log"
)

// A Leader is a privileaged node. It was appointed as a leader in an election
// and acts as a cache at the margin of the autonomous system.
// A Leader can:
//  - download from anybody
//  - upload to anybody(see race condition)
// A Leader can not:
//  - upload to different AS via an indirect connection
//
// The Leader forwards the 'have' message towards the Follower. When a Follower
// asks the Leader for a piece, the Leader with either have the piece and
// respond, or, in case it does not have it, it will try acquire the piece as
// fast as possbile. To to that, the Leader gives priority to pieces requested
// by Followers. This is achieved by using a new Picker.
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
  log.LogLeader(log.Leader{
    Id : l.Id,
  })

  l.Peer.Run(l.outgoingConnection)

  l.Picker = NewPicker(l.Storage)
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

      // We add the upload capability even though it will not be used.
      // This is necessary in case we will contacted for a bidirectional
      // leader-leader connection.

      // log.Println(l.Id, "<-", peer)
      torrent.
        NewConnector(l.Id, peer, l.Components).
        WithUpload(NewUpload).
        WithDownload(NewDownload).
        Register(l.Peer.Peer)
    }

    // Once the connections are made, we only need to register the forwarder
    l.registerForwarder(follower, peer)
  }

  // Send the messages to corresponding forwarders
  l.forward(m)

  l.Peer.RunRecv(l.GetId(m), m, l.incomingConnection)
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
    WithUpload(NewUpload).
    WithDownload(NewDownload).
    Register(l.Peer.Peer)
}

func (l *Leader) incomingConnection(id string) {
  l.outgoingConnection(id)
}
