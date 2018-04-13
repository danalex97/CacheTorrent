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

  isInterested bool // If the other peer is interested in my pieces
  choke        bool // If I choke to connection to that peer

  handshake *Handshake

  connector  *Connector
}

func NewUpload(connector *Connector) *Upload {
  return &Upload{
    Components: connector.components,

    me:        connector.from,
    to:        connector.to,
    connector: connector,

    isInterested: false, // initially, nobody is interested in my pieces
    choke:        true,  // initially, I choke all peers

    handshake: connector.handshake,
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

func (u *Upload) Choke() {
  u.choke = true
  u.Transport.ControlSend(u.to, choke{u.me})

  // Refuse to transmit
  u.handshake.uplink.Clear()
}

func (u *Upload) Unchoke() {
  u.choke = false
  u.Transport.ControlSend(u.to, unchoke{u.me})
}

func (u *Upload) interested(interested bool) {
  u.isInterested = interested

  if interested {
    u.Choker.Interested(u.connector)
  } else {
    u.Choker.NotInterested(u.connector)
  }
}
