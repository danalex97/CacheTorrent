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

func (c *Connector) WithDownloadWithRedirect(redirectId string) *Connector {
  c.Download = NewDownloadWithRedirect(c, redirectId)
  return c
}

func (c *Connector) Strip() *torrent.Connector {
  return c.Connector
}
