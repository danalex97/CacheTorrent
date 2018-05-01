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
  return p.IterateBuckers(peer, p.SelectBucket)
}

/* We use a different bucket selection which takes into account
 * the indirect requested that were made as well.
 */
func (p *CachePicker) SelectBucket(bucket map[int]bool,
                              haves map[int]bool,
                              tiebreaks map[int]int) (int, bool) {
  /*
   * @haves: set of pieces remote peer has
   * @tiebreaks: set of pieces with active started requested
   */

  look := func (check func (index int) bool) (int, bool) {
    for index, _ := range bucket {
     if !check(index) {
       continue
     }

     // if the remote peer has the piece
     if _, ok := haves[index]; ok {
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

  checkReqs := func (index int) bool {
    return p.requested[index]
  }
  checkTrue := func (index int) bool {
    return true
  }

  if ans, ok := look(checkReqs); ok {
    return ans, ok
  }
  return look(checkTrue)
}

func (p *CachePicker) IsBanned(index int) bool {
  if _, ok := p.Bans[index]; ok {
    return ok
  }

  _, ok := p.Storage.Have(index)
  if ok {
    // We also want to delete the fact that we have a request for that node
    // since is already banned.
    delete(p.requested, index)

    // we cache only positives
    p.Bans[index] = true

    // once we have a piece we can save some memory by deleting the haves of
    // those pieces
    for _, have := range p.have {
      delete(have, index)
    }
  }

  return ok
}
