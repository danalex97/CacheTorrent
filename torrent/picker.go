package torrent

/**
 This file follows the 'PiecePicker.py' file from BitTorrent 5.3.0 release.

 We follow the description in Bram Cohen's Incentives Build Robustness
 in BitTorrent, that is:
  - the policy is rarest first
  - first pieces are provided in random order rather than by rarest first policy
 */

import (
  "sync"
)

type Picker struct {
  *sync.Mutex

  have map[string]map[int]bool // the pieces that the remote peers have
}

func NewPicker(pieces []pieceMeta) *Picker {
  return &Picker{
    new(sync.Mutex),
    make(map[string]map[int]bool),
  }
}

func (p *Picker) GotHave(peer string, index int) {
  p.Lock()
  defer p.Unlock()

  if _, ok := p.have[peer]; !ok {
    p.have[peer] = make(map[int]bool)
  }
  p.have[peer][index] = true
}

func (p *Picker) Next(peer string) (index int) {
  // [TODO]
  return
}
