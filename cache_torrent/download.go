package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

type CacheDownload struct {
  *torrent.TorrentDownload
}

func NewDownload(connector *torrent.Connector) torrent.Download {
  return &CacheDownload{
    TorrentDownload : torrent.NewDownload(connector).(*torrent.TorrentDownload),
  }
}

func (u *CacheDownload) Recv(m interface {}) {
  switch m.(type) {
  default:
    u.TorrentDownload.Recv(m)
  }
}
