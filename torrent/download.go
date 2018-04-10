package torrent

import (
  // "fmt"
)

// This file follows the 'download.py' file from BitTorrent 5.3.0 release

type Download struct {
  *Components

  me   string // the node that downloads
  from string // the node that we download from

  activeRequests map[pieceMeta]bool // requests that were made, but we still
  // did not received a piece back as a response

  connector *Connector
}

func NewDownload(connector *Connector) Runner {
  return &Download{
    connector.components,

    connector.from,
    connector.to,

    make(map[pieceMeta]bool),

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
    lost := []pieceMeta{}
    for p, _ := range d.activeRequests {
      lost = append(lost, p)
    }
    d.MultiDownload.Lost(lost)

    // Since I am choked, I remove all activeRequests
    d.activeRequests = make(map[pieceMeta]bool)

    // Send interested message to node, since I am choked
    d.Transport.ControlSend(d.from, interested{d.me})
  case unchoke:
    // Request pieces from peer
    // [TODO]
  case piece:
    // Remove the request from activeRequests
    piece := pieceMeta{
      msg.index,
      msg.begin,
      msg.piece.Size,
    }
    delete(d.activeRequests, piece)

    // Request more pieces
    d.requestMore()

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

func (d *Download) requestMore() {
  //[TODO]
}
