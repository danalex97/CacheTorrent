package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

type CacheUpload struct {
  *torrent.TorrentUpload
}

func NewUpload(connector *torrent.Connector) torrent.Upload {
  return &CacheUpload{
    TorrentUpload : torrent.NewUpload(connector).(*torrent.TorrentUpload),
  }
}

func (u *CacheUpload) Recv(m interface {}) {
  switch msg := m.(type) {
  case torrent.Request:
    _, ok := u.Storage.Have(msg.Index)
    if ok {
      // We only transfer the piece if we have it.
      u.TorrentUpload.Recv(m)
    }
  default:
    u.TorrentUpload.Recv(m)
  }
}
