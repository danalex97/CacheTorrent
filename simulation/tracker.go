package simulation

import (
  . "github.com/danalex97/Speer/sdk/go"
)

type Tracker interface {
  TorrentNode

  New(id string) Tracker
}
