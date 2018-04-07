package simulation

import (
  . "github.com/danalex97/Speer/sdk/go"
)

type Peer interface {
  TorrentNode

  New(id string) Peer
}
