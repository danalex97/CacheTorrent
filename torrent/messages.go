package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
)

/* Original BitTorrent protocol messages. */
type choke struct {}
type unchoke struct {}
type interested struct {}
type notInterested struct {}

type have struct {
  index int // Index of the piece that I have
}

type request struct {
  index  int
  begin  int
  length int
}

// The actual piece as a response to a `request` message
type piece struct {
  index  int
  begin  int
  piece  Data
}

// We do not model endgame mode, so we have no `cancel` message

/* Tracker control messages. */
type join struct {
  id string
}

type neighbours struct {
  ids []string
}

type trackerReq struct {
  from string
}

type trackerRes struct {
  id string
}
