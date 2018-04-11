package torrent

/**
 This file follows the 'PiecePicker.py' file from BitTorrent 5.3.0 release.

 We follow the description in Bram Cohen's Incentives Build Robustness
 in BitTorrent, that is:
  - the policy is rarest first
  - first pieces are provided in random order rather than by rarest first policy

  Some reponsibilities of the 'RequestManager.py' have been moved to this file,
  that is the accounting of active requests.
 */

import (
  "sync"
)

type Picker struct {
  *sync.Mutex

  freq   map[int]int // map from index to frequency
  bucket map[int]map[int]bool // map from frequency to bucket
  // a bucket is a set of indexes having a specific frequency

  have   map[string]map[int]bool // the pieces that the remote peers have
  active map[int]int             // number of active requests for a piece
}

func NewPicker(pieces []pieceMeta) *Picker {
  return &Picker{
    new(sync.Mutex),
    make(map[int]int),
    make(map[int]map[int]bool),
    make(map[string]map[int]bool),
    make(map[int]int),
  }
}

func (p *Picker) GotHave(peer string, index int) {
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
  if _, ok := p.bucket[freq]; !ok {
    p.bucket[freq] = make(map[int]bool)
  }
  // erase piece from old bucket
  if _, ok := p.bucket[freq - 1]; ok {
    delete(p.bucket[freq - 1], index)
  }
  // insert piece into new bucket
  if _, ok := p.bucket[freq][index]; !ok {
    p.bucket[freq][index] = true
  }
}

func (p *Picker) Active(peer string, index int) {
  p.Lock()
  defer p.Unlock()


}

func (p *Picker) Inactive(peer string, index int) {
  p.Lock()
  defer p.Unlock()

}

func (p *Picker) Next(peer string) (index int) {
  p.Lock()
  defer p.Unlock()

  return
}
