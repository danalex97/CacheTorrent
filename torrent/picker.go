package torrent

/**
 * This file follows the 'PiecePicker.py' file from BitTorrent 5.3.0 release.
 *
 * We follow the description in Bram Cohen's Incentives Build Robustness
 * in BitTorrent, that is:
 * - the policy is rarest first
 * - first pieces are provided in random order rather than by rarest first policy
 *
 * We do not model endgame mode.
 *
 * Some reponsibilities of the 'RequestManager.py' have been moved to this file,
 * that is the accounting of active requests.
 */

import (
  "sync"
)

/* See comments below for the interface. */
type Picker interface {
  GotHave(peer string, index int)

  Active(index int)
  Inactive(index int)

  Next(peer string) (int, bool)
}

type TorrentPicker struct {
  *sync.Mutex

  Storage Storage

  freq    map[int]int // map from index to frequency
  buckets map[int]map[int]bool // map from frequency to bucket
  // a bucket is a set of indexes having a specific frequency

  have   map[string]map[int]bool // the pieces that the remote peers have
  active map[int]int             // number of active requests for a piece

  Bans   map[int]bool            // the pieces that I already have stored
}

func NewPicker(Storage Storage) Picker {
  return &TorrentPicker{
    new(sync.Mutex),
    Storage,
    make(map[int]int),
    make(map[int]map[int]bool),
    make(map[string]map[int]bool),
    make(map[int]int),
    make(map[int]bool),
  }
}

/*
 * Handler for receiving a `have` message.
 */
func (p *TorrentPicker) GotHave(peer string, index int) {
  p.Lock()
  defer p.Unlock()

  // update have
  if _, ok := p.have[peer]; !ok {
    p.have[peer] = make(map[int]bool)
  }
  p.have[peer][index] = true

  // update freq
  if _, ok := p.freq[index]; !ok {
    p.freq[index] = 0
  }
  p.freq[index] = p.freq[index] + 1

  // make bucket if not present
  freq := p.freq[index]
  if _, ok := p.buckets[freq]; !ok {
    p.buckets[freq] = make(map[int]bool)
  }
  // erase piece from old bucket
  if _, ok := p.buckets[freq - 1]; ok {
    delete(p.buckets[freq - 1], index)
  }
  // insert piece into new bucket
  if _, ok := p.buckets[freq][index]; !ok {
    p.buckets[freq][index] = true
  }
}

/*
 * Mark a certain pice as being in an active request -- that is the transfer
 * has been scheduled, but it is not yet finished.
 */
func (p *TorrentPicker) Active(index int) {
  p.Lock()
  defer p.Unlock()

  if _, ok := p.active[index]; !ok {
    p.active[index] = 0
  }
  p.active[index] = p.active[index] + 1
}

/*
 * Mark a piece as inactive -- that is the request has been eliminated or
 * the piece transfer has finished. (see `download.go`)
 */
func (p *TorrentPicker) Inactive(index int) {
  p.Lock()
  defer p.Unlock()

  p.active[index] = p.active[index] - 1
}

/*
 * Return the next piece for a certain peer.
 */
func (p *TorrentPicker) Next(peer string) (int, bool) {
  return p.IterateBuckers(peer, p.SelectBucket)
}

type Selector func (b, h map[int]bool, t map[int]int) (int, bool)

func (p *TorrentPicker) IterateBuckers(peer string, selector Selector) (int, bool) {
  p.Lock()
  defer p.Unlock()

  /*
   * @haves: set of pieces remote peer has
   * @tiebreaks: set of pieces with active started requests
   */

  haves := p.have[peer]
  if haves == nil {
    haves = make(map[int]bool)
    p.have[peer] = haves
  }
  tiebreaks := p.active

  if len(haves) == 0 {
    return 0, false
  }

  // find maximum frequency
  mx := -1
  for fr, _ := range p.buckets {
    if mx == -1 {
      mx = fr
    }

    if fr > mx {
      mx = fr
    }
  }

  // itereate through buckets from rarest to most common piece
  for fr := 1; fr <= mx; fr++ {
    bucket := p.buckets[fr]
    if bucket == nil {
      continue
    }

    index, ok := selector(bucket, haves, tiebreaks)
    if ok {
      return index, ok
    }
  }

  /**
   * We do not request a piece that was already requested. This sould not
   * increase the download time significantly assuming a small request queue
   * size.
   *
   * To fully eliminate this effect we can use config.Config.Backlog = 1.
   */

  return 0, false
}

func (p *TorrentPicker) SelectBucket(bucket map[int]bool,
                              haves map[int]bool,
                              tiebreaks map[int]int) (int, bool) {

  /*
   * @haves: set of pieces remote peer has
   * @tiebreaks: set of pieces with active started requests
   */

  for index, _ := range bucket {
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

func (p *TorrentPicker) IsBanned(index int) bool {
  if _, ok := p.Bans[index]; ok {
    return ok
  }

  _, ok := p.Storage.Have(index)
  if ok {
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
