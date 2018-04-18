package torrent

// This file follows the 'Upload.py' file from BitTorrent 5.3.0 release

import (
  . "github.com/danalex97/Speer/interfaces"
  "strconv"
)

type Upload interface {
  Runner

  Choke()    // Actions done when I choke a connection(upload)
  Unchoke()  // Actions done when I unchoke a connection(upload)

  Choking()      bool     // Returns if I'm choking the connection
  IsInterested() bool     // Returns if the other peer is interested in my pieces
  Rate()         float64  //
}

type upload struct {
  *Components

  me string
  to string

  isInterested bool // If the other peer is interested in my pieces
  choke        bool // If I choke to connection to that peer

  handshake Handshake
}

func NewUpload(connector *Connector) Upload {
  return &upload{
    Components: connector.Components,

    me:        connector.from,
    to:        connector.to,

    isInterested: false, // initially, nobody is interested in my pieces
    choke:        true,  // initially, I choke all peers

    handshake: connector.handshake,
  }
}

func (u *upload) Run() {
}

func (u *upload) Recv(m interface {}) {
  switch msg := m.(type) {
  case NotInterested:
    u.interested(false)
  case Interested:
    u.interested(true)
  case Request:
    meta ,_ := u.Storage.Have(msg.Index)

    toUpload := Data{
      strconv.Itoa(meta.Index),
      meta.Length,
    }

    // When we receive a request we can upload the piece.
    u.handshake.Uplink().Upload(toUpload)
  }
}

/*
 * Function called when we want to choke the upload connection.
 */
func (u *upload) Choke() {
  u.choke = true
  // Let the other node know
  u.Transport.ControlSend(u.to, Choke{u.me})

  // Refuse to transmit
  u.handshake.Uplink().Clear()
}

/*
 * Function called when we want to unchoke an upload.
 */
func (u *upload) Unchoke() {
  u.choke = false

  // Let the other node know
  u.Transport.ControlSend(u.to, Unchoke{u.me})
}

func (u *upload) interested(interested bool) {
  u.isInterested = interested

  if interested {
    u.Choker.Interested(u)
  } else {
    u.Choker.NotInterested(u)
  }
}

/*
 * Return if I am choking the connection.
 */
func (u *upload) Choking() bool {
  return u.choke
}

/*
 * Return if the other peer is interested in my pieces.
 */
func (u *upload) IsInterested() bool {
  return u.isInterested
}

/*
 * Returns the downoad rate of the connection.
 */
func (u *upload) Rate() float64 {
  // [TODO]
  return 0
}
