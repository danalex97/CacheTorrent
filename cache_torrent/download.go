package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
  "fmt"
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
  case Miss:
    fmt.Println("Miss")
  default:
    u.TorrentDownload.Recv(m)
  }
}
