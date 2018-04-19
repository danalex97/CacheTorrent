package cache_torrent

/**
 * Indirect requests extension.
 */

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
