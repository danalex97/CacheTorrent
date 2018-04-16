package torrent

import (
  "github.com/danalex97/nfsTorrent/config"
  "sync"
  "fmt"
)

const pieceNumber int = config.Pieces

type Storage interface {
  Have(index int) (pieceMeta, bool)
  Store(piece)

  Pieces() []int
}

type storage struct {
  sync.RWMutex

  id        string
  pieces    map[int]pieceMeta // the pieces that I have
  completed bool
}

func NewStorage(id string, pieces []pieceMeta) Storage {
  storage := new(storage)

  storage.id = id
  storage.completed = false

  storage.pieces = make(map[int]pieceMeta)
  for _, p := range pieces {
    storage.pieces[p.index] = p
  }

  storage.checkCompleted()

  return storage
}

/*
 * Returns if I have a piece(I have it downloaded and stored).
 */
func (s *storage) Have(index int) (pieceMeta, bool) {
  s.RLock()
  defer s.RUnlock()

  p, ok := s.pieces[index]
  return p, ok
}

/*
 * Store a piece from a `piece` message.
 */
func (s *storage) Store(p piece) {
  s.Lock()
  defer s.Unlock()

  s.pieces[p.index] = pieceMeta{
    p.index,
    p.begin,
    p.piece.Size,
  }

  s.checkCompleted()
}

/*
 * Return a list of all indexes of the pieces currently stored.
 */
func (s *storage) Pieces() []int {
  s.RLock()
  defer s.RUnlock()

  pieces := []int{}
  for _, piece := range s.pieces {
    pieces = append(pieces, piece.index)
  }

  return pieces
}

func (s *storage) checkCompleted() {
  if len(s.pieces) == pieceNumber && !s.completed {
    fmt.Println(s.id, "Completed")
    s.completed = true
  }
}
