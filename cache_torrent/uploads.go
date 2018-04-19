package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

type upload struct {
  torrent.Upload
}

func NewLocalUpload(c *Connector) torrent.Upload {
  return &upload{
    Upload : torrent.NewUpload(c.Connector),
  }
}
