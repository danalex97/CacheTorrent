package torrent

// This file follows the 'Upload.py' file from BitTorrent 5.3.0 release

import (
  . "github.com/danalex97/Speer/interfaces"
  "strconv"
)

type Upload struct {
  *Components

  me string
  to string

  link Link

  connector  *Connector
}

func NewUpload(connector *Connector) *Upload {
  return &Upload{
    connector.components,
    connector.from,
    connector.to,
    connector.link,
    connector,
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
    meta ,_ := u.Storage.Have(msg.index)

    toUpload := Data{
      strconv.Itoa(meta.index),
      meta.length,
    }

    // When we receive a request we can upload the piece.
    u.link.Upload(toUpload)
  }
}

func (u *Upload) interested(interested bool) {
  if interested {
    u.Choker.Interested(u.connector)
  } else {
    u.Choker.NotInterested(u.connector)
  }
}
