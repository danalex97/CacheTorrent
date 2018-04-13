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

  handshake *Handshake

  connector  *Connector
}

func NewUpload(connector *Connector) *Upload {
  return &Upload{
    connector.components,
    connector.from,
    connector.to,
    connector.handshake,
    connector,
  }
}

func (u *Upload) Run() {
}

func (u *Upload) Recv(m interface {}) {
  switch msg := m.(type) {
  case notInterested:
    u.interested(false)
  case interested:
    u.interested(true)
  case request:
    meta ,_ := u.Storage.Have(msg.index)

    toUpload := Data{
      strconv.Itoa(meta.index),
      meta.length,
    }

    // When we receive a request we can upload the piece.
    u.handshake.Uplink().Upload(toUpload)
  }
}

func (u *Upload) interested(interested bool) {
  u.connector.isInterested = interested

  if interested {
    u.Choker.Interested(u.connector)
  } else {
    u.Choker.NotInterested(u.connector)
  }
}
