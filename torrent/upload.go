package torrent

// This file follows the 'Upload' BitTorrent 5.3.0 release

type Upload struct {
  me string
  to string
}

func NewUpload(me, to string) Runner {
  upload := new(Upload)

  upload.me = me
  upload.to = to

  return upload
}

func (u *Upload) Run() {

}

func (u *Upload) Recv(m interface {}) {

}
