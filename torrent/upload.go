package torrent

// This file follows the 'Upload' BitTorrent 5.3.0 release

type Upload struct {
}

func NewUpload() Runner {
  download := new(Upload)
  return download
}

func (d *Upload) Run() {

}

func (d *Upload) Recv(m interface {}) {

}
