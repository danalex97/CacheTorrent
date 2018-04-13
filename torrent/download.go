package torrent

// This file follows the 'download.py' file from BitTorrent 5.3.0 release.

import (
  "github.com/danalex97/nfsTorrent/config"
  "strconv"
  // "runtime"
  // "fmt"
)

const backlog int = config.Backlog

type Download struct {
  *Components

  me   string // the node that downloads
  from string // the node that we download from

  activeRequests map[int]bool // requests that were made, but we still
  // did not received a piece back as a response
  handshake *Handshake

  connector *Connector
}

func NewDownload(connector *Connector) *Download {
  return &Download{
    connector.components,

    connector.from,
    connector.to,

    make(map[int]bool),
    connector.handshake,

    connector,
  }
}

func (d *Download) Run() {
  // Watch the link to deliver the piece messages
  for {
    data, ok := <-d.handshake.Downlink().Download()
    if !ok {
      // channel closed
      break
    }

    index, _ := strconv.Atoi(data.Id)
    length   := data.Size
    // assumes equal sized pieces
    begin    := index * length

    piece := piece{
      d.from,
      index,
      begin,
      data,
    }

    // send message to myself
    d.Transport.ControlSend(d.me, piece)
  }
}

func (d *Download) Recv(m interface {}) {
  switch msg := m.(type) {
  case choke:
    // Make connection as choked
    d.connector.choked = true

    // Request queued pieces that were lost from the peer that choked us
    for p, _ := range d.activeRequests {
      // let picker know
      d.Picker.Inactive(p)
    }
    // Redistribute the requests for lost pieces
    d.Choker.Lost()
    // Stop the download link as well
    // [TODO] d.handshake.Downlink().Clear()

    // Handle control messages
    if len(d.activeRequests) > 0 {
      // Send interested message to node, since I am choked
      d.interested(true)
    } else {
      // If there is no piece that I am interested in, then I am not
      // interested any more.
      _, ok := d.Picker.Next(d.from)
      d.interested(ok)
    }

    // Since I am choked, I remove all activeRequests
    d.activeRequests = make(map[int]bool)

  case unchoke:
    // Request pieces from peer
    d.connector.choked = false

    d.RequestMore()
  case piece:
    // Remove the request from activeRequests
    index := msg.index
    delete(d.activeRequests, index)

    // Let Picker know active requests changed
    d.Picker.Inactive(index)

    // Store the piece
    d.Storage.Store(msg)

    // Let the others know I have the piece
    d.Choker.Have(index)

    // We need to request more only after we stored the piece, so we don't
    // request the same thing twice.
    d.RequestMore()
  case have:
    index := msg.index

    // send interested if I'm not interested and chocked
    if d.connector.choked && !d.connector.interested {
      // I need to be interested in the piece as well
      if _, ok := d.Storage.Have(index); !ok {
        // Send interested message to node
        d.interested(true)
      }
    }

    // let picker know I can get piece index
    d.Picker.GotHave(d.from, index)
  }
}

func (d *Download) RequestMore() {
  size := backlog
  if len(d.activeRequests) >= size {
    return
  }

  // Request more pieces
  for len(d.activeRequests) < size {
    interest, ok := d.Picker.Next(d.from)
    if !ok {
      // We can't find any useful piece to request
      break
    }

    // If I'm not interested, become interested
    d.interested(true)
    d.Transport.ControlSend(d.from, request{d.me, interest})

    // Update active requests: since our network model is assumed
    // to be perfect, we assume the requests that we make to be active.
    // [see new_request @ RequestManager.py]
    d.activeRequests[interest] = true

    // Let Picker know active requests changed
    d.Picker.Active(interest)
  }
}

func (d *Download) interested(now bool) {
  before := d.connector.interested
  d.connector.interested = now

  if before != now {
    if now == true {
      d.Transport.ControlSend(d.from, interested{d.me})
    } else {
      d.Transport.ControlSend(d.from, notInterested{d.me})
    }
  }
}
