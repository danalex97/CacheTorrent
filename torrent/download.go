package torrent

// This file follows the 'download.py' file from BitTorrent 5.3.0 release.

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "strconv"
)

const backlog int = config.Backlog

type Download struct {
  *Components

  me   string // the node that downloads
  from string // the node that we download from

  interested bool // if I am interested in uploader's pieces
  choked     bool // if the peer that uploads to me chokes me

  activeRequests map[int]bool // requests that were made, but we still
  // did not received a piece back as a response
  handshake *Handshake
}

func NewDownload(connector *Connector) *Download {
  return &Download{
    Components: connector.components,

    me:   connector.from,
    from: connector.to,

    interested: false, // I am not interested in anything
    choked:     true,  // everybody chokes us

    activeRequests: make(map[int]bool),
    handshake: connector.handshake,
  }
}

func (d *Download) Run() {
  // Watch the link to deliver the piece messages
  for {
    data := <-d.handshake.Downlink().Download()
    piece := pieceFromDownload(d.from, data)

    // send message to myself (to avoid races)
    d.Transport.ControlSend(d.me, piece)
  }
}

func (d *Download) handlePending() {
  if !d.handshake.Done() {
    return
  }

  // Handle pending downloads. We deliver pieces directly
  // as we are in the Recv goroutine.
  for {
    select {
    case data := <-d.handshake.Downlink().Download():
      piece := pieceFromDownload(d.from, data)
      d.gotPiece(piece)
    default:
      return
    }
  }
}

func (d *Download) Recv(m interface {}) {
  switch msg := m.(type) {
  case choke:
    d.gotChoke(msg)
  case unchoke:
    d.gotUnchoke(msg)
  case piece:
    d.gotPiece(msg)
  case have:
    d.gotHave(msg)
  }
}

func (d *Download) gotChoke(msg choke) {
  // Handle all pending downloads
  d.handlePending()

  // Make connection as choked
  d.choked = true

  // Request queued pieces that were lost from the peer that choked us
  for p, _ := range d.activeRequests {
    // let picker know
    d.Picker.Inactive(p)
  }
  // Redistribute the requests for lost pieces
  d.Lost()

  // Handle control messages
  if len(d.activeRequests) > 0 {
    // Send interested message to node, since I am choked
    d.changeInterest(true)
  } else {
    // If there is no piece that I am interested in, then I am not
    // interested any more.
    _, ok := d.Picker.Next(d.from)
    d.changeInterest(ok)
  }

  // Since I am choked, I remove all activeRequests
  d.activeRequests = make(map[int]bool)
}

func (d *Download) gotUnchoke(msg unchoke) {
  // Request pieces from peer
  d.choked = false

  d.RequestMore()
}

func (d *Download) gotPiece(msg piece) {
  // Remove the request from activeRequests
  index := msg.index
  delete(d.activeRequests, index)

  // Let Picker know active requests changed
  d.Picker.Inactive(index)

  // Store the piece
  d.Storage.Store(msg)

  // Let the others know I have the piece
  d.Have(index)

  // We need to request more only after we stored the piece, so we don't
  // request the same thing twice.
  d.RequestMore()
}

func (d *Download) gotHave(msg have) {
  index := msg.index

  // send interested if I'm not interested and chocked
  if d.choked && !d.interested {
    // I need to be interested in the piece as well
    if _, ok := d.Storage.Have(index); !ok {
      // Send interested message to node
      d.changeInterest(true)
    }
  }

  // let picker know I can get piece index
  d.Picker.GotHave(d.from, index)
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
    d.changeInterest(true)
    d.Transport.ControlSend(d.from, request{d.me, interest})

    // Update active requests: since our network model is assumed
    // to be perfect, we assume the requests that we make to be active.
    // [see new_request @ RequestManager.py]
    d.activeRequests[interest] = true

    // Let Picker know active requests changed
    d.Picker.Active(interest)
  }
}

func (d *Download) changeInterest(now bool) {
  before := d.interested
  d.interested = now

  if before != now {
    if now == true {
      d.Transport.ControlSend(d.from, interested{d.me})
    } else {
      d.Transport.ControlSend(d.from, notInterested{d.me})
    }
  }
}

func pieceFromDownload(from string, data Data) piece {
  index, _ := strconv.Atoi(data.Id)
  length   := data.Size
  // assumes equal sized pieces
  begin    := index * length

  return piece{
    from,
    index,
    begin,
    data,
  }
}

/**
 * We moved some of the responsibility in 'MultiDownload.py',
 * 'download.py' and 'RequestManager.py' in the downloader as we only
 * need a struct which references the list of connections.
 */
func (d *Download) Lost() {
  for _, conn := range d.Manager.Downloads() {
    // We try to request more pieces only if the connection is not choked
    if !conn.choked {
      conn.RequestMore()
    }
  }
}

func (d *Download) Have(index int) {
  for _, conn := range d.Manager.Downloads() {
    d.Transport.ControlSend(conn.from, have{conn.me, index})
  }
}
