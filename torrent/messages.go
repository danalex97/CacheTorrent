package torrent

/* Original BitTorrent protocol messages. */

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
