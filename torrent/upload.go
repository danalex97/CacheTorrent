package torrent

// This file follows the 'Upload' BitTorrent 5.3.0 release

type Upload struct {
}

func NewUpload() Runner {
  upload := new(Upload)
  return upload
}

func (u *Upload) Run() {

}

func (u *Upload) Recv(m interface {}) {

}
