package torrent

// This file follows the 'Upload' BitTorrent 5.3.0 release

import (
  "fmt"
)

type Connector struct {
  from string
  to   string

  interested bool
  choked     bool

  upload   Runner
  download Runner

  components *Components
}

func NewConnector(from, to string, components *Components) Runner {
  connector := new(Connector)

  connector.from = from
  connector.to = to

  connector.components = components
  connector.upload     = NewUpload(connector)
  connector.download   = NewDownload(connector)

  connector.interested = false
  connector.choked     = true

  return connector
}

func (c *Connector) Run() {
  fmt.Println(c)

  go c.upload.Run()
  go c.download.Run()
}

func (c *Connector) Recv(m interface {}) {
  c.upload.Recv(m)
  c.download.Recv(m)
}
