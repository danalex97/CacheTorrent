package torrent

/* Original BitTorrent protocol messages. */

/* Tracker control messages. */
type join struct {
  id string
}

type neighbours struct {
  ids []string
}
