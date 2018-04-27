package torrent

/**
 * This file follows the 'download.py' file from BitTorrent 5.3.0 release.
 *
 * We moved here some of the responsibility of 'MultiDownload.py'
 * and 'RequestManager.py'.
 */

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"
  "github.com/danalex97/nfsTorrent/log"

  "container/list"
  "strconv"
)

var backlog  config.Const = config.NewConst(config.Backlog)
var removeInterval int    = 20000

type Download interface {
  Runner

  Choked() bool     // Returns if the peer that uploads to me chokes me.
  Interested() bool // Returns if I'm interested in the uploader's piece.

  Me()   string // The ID of the peer that I download from.
  From() string // The ID of the peer that I download from.

  RequestMore() // Request more pieces from a the peer.

  Rate() float64 //
}

type TorrentDownload struct {
  *Components

  me   string // the node that downloads
  from string // the node that we download from

  interested bool // if I am interested in uploader's pieces
  choked     bool // if the peer that uploads to me chokes me

  ActiveRequests map[int]bool
  /*
   * Requests that were made, but we still did not received a piece
   * back as a response.
   */

  handshake Handshake

  times *list.List
}

func NewDownload(connector *Connector) Download {
  return &TorrentDownload{
    Components: connector.Components,

    me:   connector.From,
    from: connector.To,

    interested: false, // I am not interested in anything
    choked:     true,  // everybody The ID of the peer that I download from.chokes us

    ActiveRequests: make(map[int]bool),
    handshake: connector.Handshake,

    times: list.New(),
  }
}

/*
 * Returns if the peer that uploads to me chokes me.
 */
func (d *TorrentDownload) Choked() bool {
  return d.choked
}

/*
 *  Returns if I'm interested in the uploader's piece.
 */
func (d *TorrentDownload) Interested() bool {
  return d.interested
}

/*
 * The ID of the peer that downloads.
 */
func (d *TorrentDownload) Me() string {
  return d.me
}

/*
 * The ID of the peer that I download from.
 */
func (d *TorrentDownload) From() string {
  return d.from
}

func (d *TorrentDownload) Run() {
  // Watch the link to deliver the piece messages
  for {
    data := <-d.handshake.Downlink().Download()
    piece := pieceFromDownload(d.from, data)

    // send message to myself (to avoid races)
    d.Transport.ControlSend(d.me, piece)
  }
}

func (d *TorrentDownload) handlePending() {
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

func (d *TorrentDownload) Recv(m interface {}) {
  switch msg := m.(type) {
  case Choke:
    d.gotChoke(msg)
  case Unchoke:
    d.gotUnchoke(msg)
  case Piece:
    d.gotPiece(msg)
  case Have:
    d.gotHave(msg)
  }
}

func (d *TorrentDownload) gotChoke(msg Choke) {
  // Handle all pending downloads
  d.handlePending()

  // Make connection as choked
  d.choked = true

  // Request queued pieces that were lost from the peer that choked us
  for p, _ := range d.ActiveRequests {
    // let picker know
    d.Picker.Inactive(p)
  }
  // Redistribute the requests for lost pieces
  d.lost()

  // Handle control messages
  if len(d.ActiveRequests) > 0 {
    // Send interested message to node, since I am choked
    d.changeInterest(true)
  } else {
    // If there is no piece that I am interested in, then I am not
    // interested any more.
    _, ok := d.Picker.Next(d.from)
    d.changeInterest(ok)
  }

  // Since I am choked, I remove all ActiveRequests
  d.ActiveRequests = make(map[int]bool)
}

func (d *TorrentDownload) gotUnchoke(msg Unchoke) {
  // Request pieces from peer
  d.choked = false

  d.RequestMore()
}

func (d *TorrentDownload) gotPiece(msg Piece) {
  // Update rate.
  d.updateRate()

  // Log the piece
  log.LogTransfer(log.Transfer{
    From  : d.From(),
    To    : d.Me(),
    Index : msg.Index,
  })

  // Remove the request from ActiveRequests
  index := msg.Index
  delete(d.ActiveRequests, index)

  // Let Picker know active requests changed
  d.Picker.Inactive(index)

  // Store the piece
  d.Storage.Store(msg)

  // Let the others know I have the piece
  d.have(index)

  // We need to request more only after we stored the piece, so we don't
  // request the same thing twice.
  if !d.Choked() {
    // Since the Piece control message is delivered asynchronously with the
    // Download, it may be that we are already Choked and, thus, we don't
    // need to request more pieces.
    d.RequestMore()
  }
}

func (d *TorrentDownload) gotHave(msg Have) {
  index := msg.Index

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

/*
 * Request more pieces from a the peer.
 *
 * The pieces are chosen using the Picker.
 */
func (d *TorrentDownload) RequestMore() {
  size := backlog.Int()
  if len(d.ActiveRequests) >= size {
    return
  }

  // Request more pieces
  for len(d.ActiveRequests) < size {
    interest, ok := d.Picker.Next(d.from)
    if !ok {
      // We can't find any useful piece to request
      if len(d.ActiveRequests) == 0 {
        // If we can't find any useful piece and the length of active requests
        // is 0, then we are no longer interested.
        d.changeInterest(false)
      }
      break
    }

    // If I'm not interested, become interested
    d.changeInterest(true)
    d.Transport.ControlSend(d.from, Request{d.me, interest})

    // Update active requests: since our network model is assumed
    // to be perfect, we assume the requests that we make to be active.
    // [see new_request @ RequestManager.py]
    d.ActiveRequests[interest] = true

    // Let Picker know active requests changed
    d.Picker.Active(interest)
  }
}

func (d *TorrentDownload) changeInterest(now bool) {
  before := d.interested
  d.interested = now

  if before != now {
    if now == true {
      d.Transport.ControlSend(d.from, Interested{d.me})
    } else {
      d.Transport.ControlSend(d.from, NotInterested{d.me})
    }
  }
}

func pieceFromDownload(from string, data Data) Piece {
  index, _ := strconv.Atoi(data.Id)
  length   := data.Size
  // assumes equal sized pieces
  begin    := index * length

  return Piece{
    from,
    index,
    begin,
    data,
  }
}

/*
 * Calculate the download rate as a moving average.
 */
func (d *TorrentDownload) Rate() float64 {
  return float64(d.times.Len())
}

func (d *TorrentDownload) updateRate() {
  t := d.Time()

  d.times.PushBack(t)
  for d.times.Front().Value.(int) <= t - removeInterval {
    d.times.Remove(d.times.Front())
  }
}

/**
 * We moved some of the responsibility in 'MultiDownload.py',
 * 'download.py' and 'RequestManager.py' in the downloader as we only
 * need a struct which references the list of connections.
 */
func (d *TorrentDownload) lost() {
  for _, conn := range d.Manager.Downloads() {
    // We try to request more pieces only if the connection is not choked
    if !conn.Choked() {
      conn.RequestMore()
    }
  }
}

func (d *TorrentDownload) have(index int) {
  for _, conn := range d.Manager.Uploads() {
    d.Transport.ControlSend(conn.To(), Have{conn.Me(), index})
  }
}
