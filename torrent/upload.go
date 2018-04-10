package torrent

// This file follows the 'Upload' BitTorrent 5.3.0 release

type Upload struct {
  me string
  to string

  connector *Connector
}

func NewUpload(connector *Connector) Runner {
  return &Upload{
    connector.from,
    connector.to,
    connector,
  }
}

func (u *Upload) Run() {

}

func (u *Upload) Recv(m interface {}) {

}
