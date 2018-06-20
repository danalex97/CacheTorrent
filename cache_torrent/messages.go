package cache_torrent

// Message type used by the indirect requests extension. When a Follower wants
// to start an inidirect connection via a Leader it sends a LeaderStart message.
type LeaderStart struct {
  Id   string // the ID of the local node
  Dest string // the ID of the destination
}

// Message type used by the location awareness extension. The Neighbours message
// is similar to the Neighbours message produces by BitTorrent.
type Neighbours struct {
  Ids []string
}

// Message type used by the location awareness extension. A Peer sends a
// Candidate message with its ID, upload and download rates to be propose itself
// as a leader in an election.
type Candidate struct {
  Id string

  Up   int
  Down int
}

// Message type used by the location awareness extension. The Tracker send the
// Leaders message back to Peers, letting them know who was elected leader.
type Leaders struct {
  Ids []string
}
