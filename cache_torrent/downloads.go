package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/torrent"
)

type download struct {
  torrent.Download

  redirects  []string
  transport  Transport
}

func NewDownloadWithRedirect(c *Connector) *download {
  return &download{
    Download : torrent.NewDownload(c.Connector),

    redirects : []string{},
    transport : c.Transport,
  }
}

func (d *download) AddRedirect(id string) *download {
  d.redirects = append(d.redirects, id)
  return d
}

func (d *download) Recv(m interface {}) {
  switch m.(type) {
  case torrent.Have:
    for _, id := range d.redirects {
      d.transport.ControlSend(id, m)
    }
  }

  d.Download.Recv(m)
}
