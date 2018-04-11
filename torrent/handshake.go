package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
)

type Handshake struct {
  *Components

  from  string
  to    string
}

func NewHandshake(connector *Connector) *Handshake {
  return &Handshake{
    connector.components,
    connector.from,
    connector.to,
  }
}

func (h *Handshake) Link() Link {
  return (Link)(nil)
}
