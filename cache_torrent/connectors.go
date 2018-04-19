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

func (c *Connector) WithIndirectDownload() *Connector {
  c.Download = NewIndirectDownload(c)
  return c
}
