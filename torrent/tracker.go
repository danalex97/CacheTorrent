package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "fmt"
)

const minNodes       int = config.MinNodes
const peerNeighbours int = config.PeerNeighbours

type Tracker struct {
  ids   []string
  limit int
  neigh int
}

func (t *Tracker) New(TorrentNodeUtil) TorrentNode {
  tracker := new(Tracker)

  tracker.ids   = []string{}
  tracker.limit = minNodes
  tracker.neigh = peerNeighbours

  return tracker
}

func (t *Tracker) OnJoin() {
  fmt.Println("Tracker")
}

func (t *Tracker) OnLeave() {
}
