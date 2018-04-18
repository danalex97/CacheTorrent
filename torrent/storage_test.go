package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "sort"
  "testing"
)

/* Tests. */
func TestStorage(t *testing.T) {
  s := NewStorage("", []pieceMeta{})

  s.Store(Piece{"", 0, 0, Data{"0", 10}})
  val, _ := s.Have(0)
  assertEqual(t, val, pieceMeta{0, 0, 10})

  s.Store(Piece{"", 1, 0, Data{"1", 10}})
  val, _ = s.Have(1)
  assertEqual(t, val, pieceMeta{1, 0, 10})

  s.Store(Piece{"", 2, 0, Data{"2", 10}})
  val, _ = s.Have(2)
  assertEqual(t, val, pieceMeta{2, 0, 10})

  _, ok := s.Have(3)
  assertEqual(t, ok, false)

  ps := s.Pieces()
  sort.Ints(ps)
  for k, v := range ps {
    assertEqual(t, k, v)
  }
}

func TestStorageConcurrent(t *testing.T) {
  for i := 0; i < 10; i++ {
    s := NewStorage("", []pieceMeta{})

    done := make(chan bool)
    store := func (idx int) {
      s.Store(Piece{"", idx, 0, Data{"", 10}})
      done <- true
    }

    for j := 0; j < 10; j++ {
      go store(j)
    }
    for j := 0; j < 10; j++ {
      <- done
    }

    for j := 0; j < 5; j++ {
      go func() {
        ps := s.Pieces()
        sort.Ints(ps)
        for k, v := range ps {
          assertEqual(t, k, v)
        }
      }()
    }

    for j := 0; j < 10; j++ {
      go func(idx int) {
        val, _ := s.Have(idx)
        assertEqual(t, val, pieceMeta{idx, 0, 10})
      }(j)
    }
  }
}
