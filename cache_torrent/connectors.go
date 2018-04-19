package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

type Connector struct {
  *torrent.Connector
}

func Extend(c *torrent.Connector) *Connector {
  return &Connector{
    Connector : c,
  }
}

func (c *Connector) WithLocalUpload() *Connector {
  c.Upload = NewLocalUpload(c)
  return c
}

func (c *Connector) WithLocalDownload() *Connector {
  c.Download = NewLocalDownload(c)
  return c
}
