package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
)

/* Original BitTorrent protocol messages. */
type Choke struct { id string }
type Unchoke struct { id string }
type Interested struct { id string }
type NotInterested struct { id string }

// A have message is sent to neighbours whenever a peer
// obtains a new piece.
type Have struct {
  id    string
  index int // Index of the piece that I have
}

type Request struct {
  id string

  index  int
  // We only use the index for requests instead of specifying
  // the begin and length of the data.
}

// The actual piece as a response to a `request` message
type Piece struct {
  id string

  index  int
  begin  int
  piece  Data
}

// We do not model endgame mode, so we have no `cancel` message

/* Tracker control messages. */
type Join struct {
  id string
}

type Neighbours struct {
  ids []string
}

type TrackerReq struct {
  from string
}

type TrackerRes struct {
  id string
}

/* Used to model seeds:
    - each peers sends a seed request
    - in seed response finds if its a seed how many pieces it has
 */
type SeedReq struct {
  from string
}

type SeedRes struct {
  pieces []pieceMeta
}

/* Connections */

type ConnReq struct {
  id   string
  link Link
}
