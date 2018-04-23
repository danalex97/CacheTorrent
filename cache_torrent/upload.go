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
      u.TorrentUpload.Recv(m)
      return
    }

    // If we don't have the respective message, we should choke the sender
    // so that the sender has to re-request the pieces.
    u.TorrentUpload.Choke()
  default:
    u.TorrentUpload.Recv(m)
  }
}