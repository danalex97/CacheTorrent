package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
)

/* Original BitTorrent protocol messages. */
type Choke struct { Id string }
type Unchoke struct { Id string }
type Interested struct { Id string }
type NotInterested struct { Id string }

// A have message is sent to neighbours whenever a peer
// obtains a new piece.
type Have struct {
  Id    string
  Index int // Index of the piece that I have
}

type Request struct {
  Id string

  Index  int
  // We only use the index for requests instead of specifying
  // the begin and length of the data.
}

// The actual piece as a response to a `request` message
type Piece struct {
  Id string

  Index  int
  Begin  int
  Piece  Data
}

// We do not model endgame mode, so we have no `cancel` message

/* Tracker control messages. */
type Join struct {
  Id string
}

type Neighbours struct {
  Ids []string
}

type TrackerReq struct {
  From string
}

type TrackerRes struct {
  Id string
}

/* Used to model seeds:
    - each peers sends a seed request
    - in seed response finds if its a seed how many pieces it has
 */
type SeedReq struct {
  From string
}

type SeedRes struct {
  Pieces []PieceMeta
}

/* Connections */

type ConnReq struct {
  Id   string
  Link Link
}
