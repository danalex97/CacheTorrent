package torrent

import (
  "github.com/danalex97/Speer/interfaces"
)

type pieceMeta struct {
  index  int
  begin  int
  length int
}

type Components struct {
  Picker        Picker
  Storage       Storage
  Transport     interfaces.Transport
  Choker        Choker
  Manager       Manager
}
