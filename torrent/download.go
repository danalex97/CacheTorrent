package torrent

// This file follows the 'download' BitTorrent 5.3.0 release

type Download struct {
}

func NewDownload() Runner {
  download := new(Download)
  return download
}

func (d *Download) Run() {

}

func (d *Download) Recv(m interface {}) {

}
