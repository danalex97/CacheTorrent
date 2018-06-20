package torrent

import (
  "github.com/danalex97/Speer/interfaces"
)

// PieceMeta is the metadata associated with a piece.
type PieceMeta struct {
  Index  int
  Begin  int
  Length int
}

// The Components structure is a container keeping all global data
// structures associated with a Peer.
type Components struct {
  // Internal data structures.
  Picker        Picker
  Storage       Storage
  Choker        Choker
  Manager       Manager

  // External(simulated) data structures.
  Transport     interfaces.Transport
  Time          func() int
}
