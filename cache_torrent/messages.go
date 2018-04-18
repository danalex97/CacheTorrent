package cache_torrent

/**
 * Used for location awareness.
 */
type LocalReq struct {
  Id string
}

type LocalRes struct {
  Ids []string
}
