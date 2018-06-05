package torrent

import (
  "github.com/danalex97/nfsTorrent/config"
  "github.com/danalex97/nfsTorrent/log"
  "sync"
  "fmt"
)

var pieceNumber config.Const = config.NewConst(config.StoragePieces)

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

  // Used for logging
  time  func() int

  // Used only for monitoring
  percents    []int
  percentDone []bool
}

func NewStorage(id string, pieces []PieceMeta, time func() int) Storage {
  storage := new(storage)

  storage.id = id
  storage.completed = false
  storage.time = time

  if len(pieces) > 0 {
    fmt.Println(id, "Pieces: ", len(pieces))
  }
  storage.pieces = make(map[int]PieceMeta)
  for _, p := range pieces {
    storage.pieces[p.Index] = p
  }

  // Used only informatively.
  storage.percents = []int{2, 20, 50, 70}
  storage.percentDone = []bool{false, false, false, false}

  storage.checkCompleted()

  return storage
}

// Returns if I have a piece(I have it downloaded and stored).
func (s *storage) Have(index int) (PieceMeta, bool) {
  s.RLock()
  defer s.RUnlock()

  p, ok := s.pieces[index]
  return p, ok
}

// Store a piece from a `piece` message.
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

// Return a list of all indexes of the pieces currently stored.
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
  // Only informative.
  for i := range s.percents {
    if s.percentDone[i] == false {
      if len(s.pieces) > pieceNumber.Int() * s.percents[i] / 100 {
        s.percentDone[i] = true
        log.Println(s.id, "Downloaded", s.percents[i], "%")
      }
    }
  }

  if len(s.pieces) == pieceNumber.Int() && !s.completed {
    // Notify logger
    time := s.time()
    log.LogCompleted(log.Completed{
      Id   : s.id,
      Time : time,
    })

    // Callback used to interact with the simulation
    config.Config.SharedCallback()

    // Notify completed
    log.Println(s.id, "Completed at", time, "ms")
    s.completed = true
  }
}
