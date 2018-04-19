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

func (c *Connector) WithDownloadWithRedirect() *Connector {
  c.Download = NewDownloadWithRedirect(c)
  return c
}

func (c *Connector) WithUploadWithRedirect() *Connector {
  c.Upload = NewUploadWithRedirect(c)
  return c
}

func (c *Connector) Strip() *torrent.Connector {
  return c.Connector
}
