package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

type download struct {
  torrent.Download
}

func NewIndirectDownload(c *Connector) torrent.Download {
  return &download{
    Download : torrent.NewDownload(c.Connector),
  }
}
