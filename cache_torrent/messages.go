package cache_torrent

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
