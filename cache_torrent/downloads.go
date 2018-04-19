package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

type download struct {
  torrent.Download
}

func NewLocalDownload(c *Connector) torrent.Download {
  return &download{
    Download : torrent.NewDownload(c.Connector),
  }
}
