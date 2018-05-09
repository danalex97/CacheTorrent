package torrent

// This file follows the 'Upload.py' file from BitTorrent 5.3.0 release

import (
  . "github.com/danalex97/Speer/interfaces"
  "strconv"
  "sync"
)

type Upload interface {
  Runner

  Choke()    // Actions done when I choke a connection(upload)
  Unchoke()  // Actions done when I unchoke a connection(upload)

  Me() string
  To() string

  Choking()      bool     // Returns if I'm choking the connection
  IsInterested() bool     // Returns if the other peer is interested in my pieces
  Rate()         float64  //
}

type TorrentUpload struct {
  *Components

  me string
  to string

  interestedMutex *sync.Mutex
  chokeMutex *sync.Mutex

  isInterested bool // If the other peer is interested in my pieces
  choke        bool // If I choke to connection to that peer

  handshake Handshake

  // Connector -- used only to reference download
  connector *Connector
}

func NewUpload(connector *Connector) Upload {
  return &TorrentUpload{
    Components: connector.Components,

    me:        connector.From,
    to:        connector.To,

    interestedMutex: new(sync.Mutex),
    chokeMutex: new(sync.Mutex),

    isInterested: false, // initially, nobody is interested in my pieces
    choke:        true,  // initially, I choke all peers

    handshake: connector.Handshake,
    connector: connector,
  }
}

func (u *TorrentUpload) Run() {
}

func (u *TorrentUpload) Recv(m interface {}) {
  switch msg := m.(type) {
  case NotInterested:
    u.interested(false)
  case Interested:
    u.interested(true)
  case Request:
    meta, _ := u.Storage.Have(msg.Index)

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
func (u *TorrentUpload) Choke() {
  u.chokeMutex.Lock()
  defer u.chokeMutex.Unlock()

  u.choke = true
  // Let the other node know
  u.Transport.ControlSend(u.to, Choke{u.me})

  // Refuse to transmit
  u.handshake.Uplink().Clear()
}

/*
 * Function called when we want to unchoke an upload.
 */
func (u *TorrentUpload) Unchoke() {
  u.chokeMutex.Lock()
  defer u.chokeMutex.Unlock()

  u.choke = false

  // Let the other node know
  u.Transport.ControlSend(u.to, Unchoke{u.me})
}

func (u *TorrentUpload) interested(interested bool) {
  u.interestedMutex.Lock()
  u.isInterested = interested
  u.interestedMutex.Unlock()

  if interested {
    u.Choker.Interested(u)
  } else {
    u.Choker.NotInterested(u)
  }
}

/*
 * The ID of the peer that uploads.
 */
func (u *TorrentUpload) Me() string {
  return u.me
}

/*
 * The ID of the peer that I upload to.
 */
func (u *TorrentUpload) To() string {
  return u.to
}

/*
 * Return if I am choking the connection.
 */
func (u *TorrentUpload) Choking() bool {
  u.chokeMutex.Lock()
  defer u.chokeMutex.Unlock()

  return u.choke
}

/*
 * Return if the other peer is interested in my pieces.
 */
func (u *TorrentUpload) IsInterested() bool {
  u.interestedMutex.Lock()
  defer u.interestedMutex.Unlock()

  return u.isInterested
}

/*
 * Returns the downoad rate of the connection.
 */
func (u *TorrentUpload) Rate() float64 {
  if u.connector.Download == nil {
    return 0
  }
  return u.connector.Download.Rate()
}
