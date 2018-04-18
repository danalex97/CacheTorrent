package simulation

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "github.com/danalex97/nfsTorrent/cache_torrent"
)

type simulatedNode struct {
}

func (s *simulatedNode) New(util TorrentNodeUtil) TorrentNode {
  if util.Join() == "" {
    return new(torrent.Tracker).New(util)
  } else {
    return new(torrent.Peer).New(util)
  }
}

type simulatedCachedNode struct {
}

func (s *simulatedCachedNode) New(util TorrentNodeUtil) TorrentNode {
  if util.Join() == "" {
    return new(torrent.Tracker).New(util)
  } else {
    return new(cache_torrent.Peer).New(util)
  }
}
