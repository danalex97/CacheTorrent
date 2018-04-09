package torrent

// This file follows the 'Upload' BitTorrent 5.3.0 release

import (
  "fmt"
)

type Connector struct {
  from string
  to   string

  upload   Runner
  download Runner
}

func NewConnector(from, to string) Runner {
  connector := new(Connector)

  connector.from = from
  connector.to = to

  connector.upload   = NewUpload()
  connector.download = NewDownload()

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
