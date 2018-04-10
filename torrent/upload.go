package torrent

// This file follows the 'Upload' file from BitTorrent 5.3.0 release

import (
  . "github.com/danalex97/Speer/interfaces"
)

type Upload struct {
  me string
  to string

  connector  *Connector
  components *Components
}

func NewUpload(connector *Connector) Runner {
  return &Upload{
    connector.from,
    connector.to,
    connector,
    connector.components,
  }
}

func (u *Upload) Run() {
}

func (u *Upload) Recv(m interface {}) {
  switch msg := m.(type) {
  case notInterested:
    u.connector.interested = false
    u.interested(u.connector.interested)
  case interested:
    u.connector.interested = true
    u.interested(u.connector.interested)
  case request:
    meta ,_ := u.components.Storage.Have(msg.index)

    toUpload := Data{
      string(meta.index),
      meta.length,
    }

    u.upload(u.to, meta, toUpload)
  }
}

func (u *Upload) interested(interested bool) {
  // Let checker know
  // [TODO]
}

func (u *Upload) upload(id string, meta pieceMeta, upload Data) {
  // [TODO]
}
