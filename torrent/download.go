package torrent

// This file follows the 'download.py' file from BitTorrent 5.3.0 release

import (
  "github.com/danalex97/nfsTorrent/config"
)

const backlog int = config.Backlog

type Download struct {
  *Components

  me   string // the node that downloads
  from string // the node that we download from

  activeRequests map[int]bool // requests that were made, but we still
  // did not received a piece back as a response

  connector *Connector
}

func NewDownload(connector *Connector) Runner {
  return &Download{
    connector.components,

    connector.from,
    connector.to,

    make(map[int]bool),

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
    lost := []int{}
    for p, _ := range d.activeRequests {
      lost = append(lost, p)
    }
    d.MultiDownload.Lost(lost)

    // Since I am choked, I remove all activeRequests
    d.activeRequests = make(map[int]bool)

    // Send interested message to node, since I am choked
    d.Transport.ControlSend(d.from, interested{d.me})
  case unchoke:
    // Request pieces from peer
    d.connector.choked = false

    d.requestMore()
  case piece:
    // Remove the request from activeRequests
    index := msg.index
    delete(d.activeRequests, index)

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
  size := backlog
  if len(d.activeRequests) >= size {
    return
  }

  // Request more pieces
  for len(d.activeRequests) < size {
    interest := d.Picker.Next(d.me)

    // If I'm not interested, become interested
    if !d.connector.interested {
      d.Transport.ControlSend(d.from, interested{d.me})
    }

    // Send request
    d.Transport.ControlSend(d.from, request{d.me, interest})

    // Update active requests: since our network model is assumed
    // to be perfect, we assume the requests that we make to be active.
    // [see new_request @ RequestManager.py]
    d.activeRequests[interest] = true
  }
}
