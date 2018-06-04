package simulation

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"

  "github.com/danalex97/nfsTorrent/torrent"
  "github.com/danalex97/nfsTorrent/cache_torrent"
  "github.com/danalex97/nfsTorrent/bias_torrent"
  "github.com/danalex97/nfsTorrent/multi_torrent"
)

type SimulatedNode struct {
}

func (s *SimulatedNode) New(util TorrentNodeUtil) TorrentNode {
  if util.Join() == "" {
    return new(torrent.Tracker).New(util)
  } else {
    config.Config.SharedInit()
    return new(torrent.Peer).New(util)
  }
}

type SimulatedBiasedNode struct {
}

func (s *SimulatedBiasedNode) New(util TorrentNodeUtil) TorrentNode {
  if util.Join() == "" {
    return new(bias_torrent.Tracker).New(util)
  } else {
    config.Config.SharedInit()
    return new(torrent.Peer).New(util)
  }
}


type SimulatedCachedNode struct {
}

func (s *SimulatedCachedNode) New(util TorrentNodeUtil) TorrentNode {
  if util.Join() == "" {
    return new(cache_torrent.Tracker).New(util)
  } else {
    config.Config.SharedInit()
    return new(cache_torrent.Peer).New(util)
  }
}

type SimulatedBiasedCachedNode struct {
}

func (s *SimulatedBiasedCachedNode) New(util TorrentNodeUtil) TorrentNode {
  if util.Join() == "" {
    return new(bias_torrent.CacheTracker).New(util)
  } else {
    config.Config.SharedInit()
    return new(cache_torrent.Peer).New(util)
  }
}

type SimulatedMultiNode struct {
}

func (s *SimulatedMultiNode) New(util TorrentNodeUtil) TorrentNode {
  if util.Join() == "" {
    // This will be replaced with the new tracker
    return new(cache_torrent.Tracker).New(util)
  } else {
    config.Config.SharedInit()
    return new(multi_torrent.MultiPeer).New(util)
  }
}
