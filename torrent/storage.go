package torrent

import (
  "github.com/danalex97/nfsTorrent/config"
  "sync"
  "fmt"
)

const pieceNumber int = config.Pieces

type Storage struct {
  sync.RWMutex

  id        string
  pieces    map[int]pieceMeta // the pieces that I have
  completed bool
}

func NewStorage(id string, pieces []pieceMeta) *Storage {
  storage := new(Storage)

  storage.id = id
  storage.completed = false

  storage.pieces = make(map[int]pieceMeta)
  for _, p := range pieces {
    storage.pieces[p.index] = p
  }

  storage.checkCompleted()

  return storage
}

func (s *Storage) Have(index int) (pieceMeta, bool) {
  s.RLock()
  defer s.RUnlock()

  p, ok := s.pieces[index]
  return p, ok
}

func (s *Storage) Store(p piece) {
  s.Lock()
  defer s.Unlock()

  s.pieces[p.index] = pieceMeta{
    p.index,
    p.begin,
    p.piece.Size,
  }

  s.checkCompleted()
}

func (s *Storage) checkCompleted() {
  if len(s.pieces) == pieceNumber && !s.completed {
    fmt.Println(s.id, "Completed")
    s.completed = true
  }
}
