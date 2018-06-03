package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/cache_torrent"
)

// A MultiPeer is wrapper over multiple Peers which follow the CacheTorrent
// protocol. The Original message IDs are decorated with an internal ID, each
// MultiPeer having an ID format "<multipeer-id>.<peer-id>". Once a message
// arrives to a Peer, the id is stripped to "<multipeer-id>".
type MultiPeer struct {
  *cache_torrent.Peer

  // Map from internal ID to Peer.
  peers map[string]cache_torrent.Peer

  // Utility structure to pass at single Peer initilization.
  util TorrentNodeUtil
}

func (p *MultiPeer) New(util TorrentNodeUtil) TorrentNode {
  return &MultiPeer{
    peers : make(map[string]cache_torrent.Peer),
    util  : util,
  }
}

func (p *MultiPeer) OnJoin() {
  if p.Transport == nil {
    return
  }

  // go func() {
  //   p.Init()
  //   go p.CheckMessages(p.Bind, p.Process)
  // }()
}
