package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

// The CacheUpload is used by Leaders. It acts like a BitTorrent upload, but
// it may be that it does not have a Piece that it advertised. In this case, it
// will not transfer the piece.
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
