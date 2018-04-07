package simulation

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
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
