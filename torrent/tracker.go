package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "fmt"
)

type Tracker struct {
}

func (t *Tracker) OnJoin() {
  fmt.Println("Tracker")
}

func (t *Tracker) OnLeave() {
}

func (t *Tracker) New(TorrentNodeUtil) TorrentNode {
  return new(Tracker)
}
