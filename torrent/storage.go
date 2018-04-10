package torrent

import (
  "sync"
)

type Storage struct {
  sync.RWMutex

  pieces map[int]pieceMeta // the pieces that I have
}

func NewStorage(pieces []pieceMeta) *Storage {
  storage := new(Storage)

  storage.pieces = make(map[int]pieceMeta)
  for _, p := range pieces {
    storage.pieces[p.index] = p
  }
  
  return storage
}

func (s *Storage) Have(index int) (pieceMeta, bool) {
  s.RLock()
  defer s.RUnlock()

  p, ok := s.pieces[index]
  return p, ok
}
