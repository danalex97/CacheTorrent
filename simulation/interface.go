package simulation

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "github.com/danalex97/nfsTorrent/cache_torrent"
)

type SimulatedNode struct {
}

func (s *SimulatedNode) New(util TorrentNodeUtil) TorrentNode {
  if util.Join() == "" {
    return new(torrent.Tracker).New(util)
  } else {
    return new(torrent.Peer).New(util)
  }
}

type SimulatedCachedNode struct {
}

func (s *SimulatedCachedNode) New(util TorrentNodeUtil) TorrentNode {
  if util.Join() == "" {
    return new(cache_torrent.Tracker).New(util)
  } else {
    return new(cache_torrent.Peer).New(util)
  }
}
