package torrent

import (
  // "fmt"
)

// This file follows the 'download.py' file from BitTorrent 5.3.0 release

type Download struct {
  *Components

  me   string // the node that downloads
  from string // the node that we download from

  connector *Connector
}

func NewDownload(connector *Connector) Runner {
  return &Download{
    connector.components,

    connector.from,
    connector.to,

    connector,
  }
}

func (d *Download) Run() {
}

func (d *Download) Recv(m interface {}) {
  switch msg := m.(type) {
  case choke:
    // Make connection as choked
    d.connector.choked = true

    // Request queued pieces that were lost from the peer that choked us
    // [TODO]

    // Send interested message to node
    d.Transport.ControlSend(d.from, interested{d.me})
  case unchoke:
    // Request pieces from peer
    // [TODO]
  case piece:
    // Store the piece
    d.Storage.Store(msg)
  case have:
    index := msg.index

    // send interested if I'm not interested and chocked
    if d.connector.choked && !d.connector.interested {
      // I need to be interested in the piece as well
      if _, ok := d.Storage.Have(index); !ok {
        // Send interested message to node
        d.connector.interested = true
        d.Transport.ControlSend(d.from, interested{d.me})
      }
    }

    // let picker know I can get piece index
    d.Picker.GotHave(d.from, index)
  }
}
