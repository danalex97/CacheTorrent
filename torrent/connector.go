package torrent

// This file follows the 'Upload' BitTorrent 5.3.0 release

import (
  "fmt"
)

type Connector struct {
}

func NewConnector() Runner {
  connector := new(Connector)
  return connector
}

func (c *Connector) Run() {
}

func (c *Connector) Recv(m interface {}) {
  fmt.Println("Connector", c)
}
