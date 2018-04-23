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
  // "fmt"
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
  // fmt.Println("Leader running.")
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

      // fmt.Println(l.Id, "<-", peer)
      torrent.
        NewConnector(l.Id, peer, l.Components).
        WithUpload(NewUpload). // [?]
        WithDownload(NewDownload).
        Register(l.Peer.Peer)
    }

    // Once the connections are made, we only need to register the forwarder
    l.registerForwarder(follower, peer)
  }

  // Send the messages to corresponding forwarders
  l.forward(m)

  // Add upload component if necessary
  // l.addUploader(m)

  l.Peer.RunRecv(m, l.incomingConnection)
}

// func (l *Leader) addUploader(m interface {}) {
//   // If we have an incoming connection, we may need to upgrade the current
//   // connection by adding a upload component.
//   id := l.GetId(m)
//   if conn, ok := l.Connectors[id]; ok && conn.(*torrent.Connector).Upload == nil {
//     // If there is a connection with the id
//     c := conn.(*torrent.Connector)
//
//     // This is ugly...
//     c.Upload = NewUpload(c)
//     go c.Upload.Run()
//     go c.Handshake.Run()
//   }
// }

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
  // fmt.Println(l.Id, "<->", id)

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
