package torrent

import (
  "github.com/danalex97/Speer/interfaces"
)

type PieceMeta struct {
  Index  int
  Begin  int
  Length int
}

type Components struct {
  Picker        Picker
  Storage       Storage
  Transport     interfaces.Transport
  Choker        Choker
  Manager       Manager
}
