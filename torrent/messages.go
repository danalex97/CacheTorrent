package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
)

// Original BitTorrent protocol message used to notify a node that it has been
// choked. See the 'Messsage' section for a list with all the types of
// messages.
type Choke struct { Id string }
// Original BitTorrent protocol message used to notify a node that it has been
// unchoked. See the 'Messsage' section for a list with all the types of
// messages.
type Unchoke struct { Id string }
// Original BitTorrent protocol message used by a node to show that it is
// interested in at least one of the messages advertised via 'Have' messages by
// another node. See the 'Messsage' section for a list with all the types of
// messages.
type Interested struct { Id string }
// Original BitTorrent protocol message used by a node to show that it is not
// interested in any of the messages advertised via 'Have' messages by another
// node. NotInterested messages are also sent after all useful Pieces have been
// received to mark download completion. See the 'Messsage' section for a list
// with all the types of messages.
type NotInterested struct { Id string }

// A 'Have' message is sent to neighbours whenever a peer obtains a new piece.
// This message is part of the Original BitTorrent protocol. See the 'Messsage'
// section for a list with all the types of messages.
type Have struct {
  Id    string  // The Id of the sender.
  Index int     // Index of the piece that the sender has.
}

// A 'Request' message is sent when a node is unchoked and a piece was
// advertised by a Peer. The node will request a piece with a specific index.
// This message is part of the Original BitTorrent protocol. See the 'Messsage'
// section for a list with all the types of messages.
type Request struct {
  Id string

  // We only use the index for requests instead of specifying
  // the begin and length of the data.
  Index  int
}

// The meta-data associated with a piece sent as a response to a
// `Request` message. This message is part of the Original BitTorrent protocol.
// See the 'Messsage' section for a list with all the types of messages.
type Piece struct {
  Id string

  // The index assoicated with the piece. Each Piece is identified by a unique
  // index as part of a file.
  Index  int
  // The index of the first byte in the file. Typically the piece sizes are
  // equal, so Begin = Index * PieceLength.
  Begin  int
  // The actual data associated with the piece. In our case, this is associated
  // with Speer simulator's interfaces.Data type.
  Piece  Data
}

// The Join message is a Tracker control message. This message is sent by a
// node to the tracker when it wants to join the torrent. See the 'Messsage'
// section for a list with all the types of messages.
type Join struct {
  Id string
}

// The Neighbours message is a Tracker control message. This message is sent
// by the tracker as a response to a Join message. See the 'Messsage' section
// for a list with all the types of messages.
type Neighbours struct {
  Ids []string
}

// The TrackerReq is a messages used for bootstraping. When a node joins the
// simulated network, it has access to only one node. The message will be sent
// to that node as question 'Who is the Tracker?'. See the 'Messsage' section
// for a list with all the types of messages.
type TrackerReq struct {
  From string
}

// The TrackerRes is a messages used for bootstraping. When a node joins the
// simulated network, it has access to only one node. The message is a response
// to a TrackerReq question 'Who is the Tracker?'. See the 'Messsage' section
// for a list with all the types of messages.
type TrackerRes struct {
  Id string
}

// The SeedReq is a message used for intergration with Speer. For simplicity,
// a Tracker will decide which nodes are seeds. This message is useless in a
// real deployment. Each peer sends a 'seed request' and in the 'seed response'
// finds if it's a seed and how many pieces it has. See the 'Messsage' section
// for a list with all the types of messages.
type SeedReq struct {
  From string
}

// The SeedReq is a message used for intergration with Speer. See 'SeedReq'
// for details. See the 'Messsage' section for a list with all the types
// of messages.
type SeedRes struct {
  Pieces []PieceMeta
}

// Connection establishment used for a three way handshake. For greater detail
// on how handshakes work, see the Handshake interface. See the 'Messsage'
// section for a list with all the types of messages.
type ConnReq struct {
  Id   string
  Link Link
}
