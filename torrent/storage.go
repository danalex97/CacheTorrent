package torrent

import (
  "github.com/danalex97/nfsTorrent/config"
  "sync"
  "fmt"
)

var pieceNumber int = config.Config.Pieces

var shared         interface {}       = config.Config.Shared
var sharedCallback func(interface {}) = config.Config.SharedCallback

type Storage interface {
  Have(index int) (PieceMeta, bool)
  Store(Piece)

  Pieces() []int
}

type storage struct {
  sync.RWMutex

  id        string
  pieces    map[int]PieceMeta // the pieces that I have
  completed bool
}

func NewStorage(id string, pieces []PieceMeta) Storage {
  storage := new(storage)

  storage.id = id
  storage.completed = false

  storage.pieces = make(map[int]PieceMeta)
  for _, p := range pieces {
    storage.pieces[p.Index] = p
  }

  storage.checkCompleted()

  return storage
}

/*
 * Returns if I have a piece(I have it downloaded and stored).
 */
func (s *storage) Have(index int) (PieceMeta, bool) {
  s.RLock()
  defer s.RUnlock()

  p, ok := s.pieces[index]
  return p, ok
}

/*
 * Store a piece from a `piece` message.
 */
func (s *storage) Store(p Piece) {
  s.Lock()
  defer s.Unlock()

  s.pieces[p.Index] = PieceMeta{
    p.Index,
    p.Begin,
    p.Piece.Size,
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
    pieces = append(pieces, piece.Index)
  }

  return pieces
}

func (s *storage) checkCompleted() {
  if len(s.pieces) == pieceNumber && !s.completed {
    // Callback used to interact with the simulation
    sharedCallback(shared)

    // Notify completed
    fmt.Println(s.id, "Completed")
    s.completed = true
  }
}
