package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
  "fmt"
)

type download struct {
  torrent.Download

  redirectId string
  transport  Transport
}

func NewDownloadWithRedirect(c *Connector, redirectId string) torrent.Download {
  return &download{
    Download : torrent.NewDownload(c.Connector),

    redirectId : redirectId,
    transport  : c.Transport,
  }
}

func (d *download) Recv(m interface {}) {
  switch m.(type) {
  case torrent.Have:
    fmt.Println("Indirect have:", m)
    d.transport.ControlSend(d.redirectId, m)
  }

  d.Download.Recv(m)
}
