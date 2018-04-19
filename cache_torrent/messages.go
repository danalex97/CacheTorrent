package cache_torrent

/**
 * Indirect requests extension.
 */

type LeaderStart struct {
  Id   string // the ID of the local node
  Dest string // the ID of the destination
}

type RemoteStart struct {
  Id string // the ID of the leader
}

type IndirectReq struct {
  From string
  Dest string

  Index int
}

/**
 * Location awareness extension.
 */

type Neighbours struct {
  Ids []string
}

type Candidate struct {
  Id string

  Up   int
  Down int
}

type Leaders struct {
  Ids []string
}
