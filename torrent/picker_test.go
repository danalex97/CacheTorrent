package torrent

import (
  "testing"
)

/* Tests. */
func TestPickerChoosesRarestPiece(t *testing.T) {
  p := NewPicker(NewStorage("", []pieceMeta{}))

  have := p.GotHave

  have("1", 5)
  have("2", 5)
  have("3", 5)
  have("1", 2)
  have("2", 2)
  have("3", 7)

  piece, ok := p.Next("3")
  assertEqual(t, piece, 7)
  assertEqual(t, ok, true)

  piece, ok = p.Next("1")
  assertEqual(t, piece, 2)
  assertEqual(t, ok, true)

  piece, ok = p.Next("2")
  assertEqual(t, piece, 2)
  assertEqual(t, ok, true)
}

func TestPickerChoosesRarestPieceConcurrently(t *testing.T) {
  for i := 0; i < 10; i++ {
    p := NewPicker(NewStorage("", []pieceMeta{}))

    done := make(chan bool)
    have := func(peer string, idx int) {
      p.GotHave(peer, idx)
      done <- true
    }

    go have("1", 5)
    go have("2", 5)
    go have("3", 5)
    go have("1", 2)
    go have("2", 2)
    go have("3", 7)

    for j := 0; j < 6; j++ {
      <- done
    }

    go func() {
      piece, ok := p.Next("3")
      assertEqual(t, piece, 7)
      assertEqual(t, ok, true)
    }()

    go func() {
      piece, ok := p.Next("1")
      assertEqual(t, piece, 2)
      assertEqual(t, ok, true)
    }()

    go func() {
      piece, ok := p.Next("2")
      assertEqual(t, piece, 2)
      assertEqual(t, ok, true)
    }()
  }
}

func TestPickerDoesntChoosePendingRequest(t *testing.T) {
  /*
   * Note: see Next specification to understand this test.
   */
  p := NewPicker(NewStorage("", []pieceMeta{}))

  have := p.GotHave

  have("1", 5)
  have("2", 5)
  have("3", 5)
  have("1", 2)
  have("2", 2)
  have("3", 7)

  p.Active(7)
  p.Active(2)
  p.Active(7)
  p.Inactive(7)

  piece, ok := p.Next("2")
  assertEqual(t, piece, 5)
  assertEqual(t, ok, true)

  p.Active(5)

  _, ok = p.Next("1")
  assertEqual(t, ok, false)

  p.Inactive(5)

  piece, ok = p.Next("1")
  assertEqual(t, piece, 5)
  assertEqual(t, ok, true)
}
