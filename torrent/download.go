package torrent

// This file follows the 'download' BitTorrent 5.3.0 release

type Download struct {
  me   string // the node that downloads
  from string // the node that we download from
}

func NewDownload(me, from string) Runner {
  download := new(Download)

  download.me   = me
  download.from = from

  return download
}

func (d *Download) Run() {

}

func (d *Download) Recv(m interface {}) {

}
