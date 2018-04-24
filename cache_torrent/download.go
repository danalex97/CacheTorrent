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

func (d *CacheDownload) Recv(m interface {}) {
  switch msg := m.(type) {

  case torrent.Have:
    index := msg.Index
    // If we have already requested the piece, it must be that the piece
    // advertisment was an indirect one.
    if _, ok := d.ActiveRequests[index]; ok {
      // Since we have active requests, it must be that we are not Choked.
      if d.Choked() {
        panic("Got active requests while Choked.")
      }

      // Since we got the second advertisment, it must be that the Leader has
      // the piece. Therefore, we must request it again for the transfer to begin.
      d.Transport.ControlSend(d.From(), torrent.Request{
        Id    : d.Me(),
        Index : index,
      })
    } else {
      // In case this is the first Have, we use the parent's handler
      d.TorrentDownload.Recv(m)
    }

  default:
    d.TorrentDownload.Recv(m)
  }
}
