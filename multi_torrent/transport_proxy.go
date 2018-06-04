package multi_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
)

type TransportProxy struct {
  Transport

  id string
}

func NewTransportProxy(t Transport, id string) *TransportProxy {
  return &TransportProxy{
    Transport : t,
    id : id,
  }
}
