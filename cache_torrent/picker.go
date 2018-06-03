package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

type CachePicker struct {
  *torrent.TorrentPicker

  requested map[int]bool // if the piece was indirectly requested
}

func NewPicker(storage torrent.Storage) torrent.Picker {
  return &CachePicker{
    TorrentPicker : torrent.NewPicker(storage).(*torrent.TorrentPicker),
    requested     : make(map[int]bool),
  }
}

func (p *CachePicker) GotRequest(index int) {
  p.Lock()
  defer p.Unlock()

  if _, ok := p.requested[index]; !ok {
    p.requested[index] = true
  }
}

func (p *CachePicker) Next(peer string) (int, bool) {
  p.RLock()
  defer p.RUnlock()

  return p.IterateBuckers(peer, p.SelectBucket)
}

// We use a different bucket selection which takes into account
// the indirect requested that were made as well.
func (p *CachePicker) SelectBucket(bucket map[int]bool,
                              haves map[int]bool,
                              tiebreaks map[int]int) (int, bool) {
  // @haves: set of pieces remote peer has
  // @tiebreaks: set of pieces with active started requested

  iterate := bucket
  check1  := haves
  check2  := p.requested
  if len(check1) < len(iterate) {
    check1, iterate = iterate, check1
  }
  if len(check2) < len(iterate) {
    check2, iterate = iterate, check2
  }
  if len(check2) < len(check1) {
    check2, check1 = check1, check2
  }

  for index, _ := range iterate {
    if _, ok := check1[index]; ok {
      if _, ok := check2[index]; ok {
        if nbr, ok := tiebreaks[index]; !ok || nbr == 0 {
          // and the piece is not already stored
          if !p.IsBanned(index) {
           return index, true
          }
        }
      }
    }
  }

  iterate  = bucket
  check   := haves
  if len(bucket) < len(haves) {
    check, iterate = iterate, check
  }
  for index, _ := range iterate {
    // if the remote peer has the piece
    if _, ok := check[index]; ok {
      // and I did not requested the piece already
      if nbr, ok := tiebreaks[index]; !ok || nbr == 0 {
        // and the piece is not already stored
        if !p.IsBanned(index) {
          return index, true
        }
      }
    }
  }

  return 0, false
}

func (p *CachePicker) IsBanned(index int) bool {
  if _, ok := p.Bans[index]; ok {
    return ok
  }

  _, ok := p.Storage.Have(index)
  if ok {
    p.RUnlock()
    p.Lock()
    defer func() {
      p.Unlock()
      p.RLock()
    }()

    // We also want to delete the fact that we have a request for that node
    // since is already banned.
    delete(p.requested, index)

    // we cache only positives
    p.Bans[index] = true

    // once we have a piece we can save some memory by deleting the haves of
    // those pieces
    for _, have := range p.Have {
      delete(have, index)
    }
  }

  return ok
}
