package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

type upload struct {
  torrent.Upload
}

func NewUploadWithRedirect(c *Connector) *upload {
  return &upload{
    Upload : torrent.NewUpload(c.Connector),
  }
}
