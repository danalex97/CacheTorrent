package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  . "github.com/danalex97/Speer/capacity"
  "testing"
)

/* Mocks. */
type mockUpload struct {
  mockRunner

  isInterested bool
  choke        bool
  rate         float64

  me string
  to string
}

func (m *mockUpload) Me() string {
  return m.me
}

func (m *mockUpload) To() string {
  return m.to
}

func (m *mockUpload) Choke() {
  m.choke = true
}

func (m *mockUpload) Unchoke() {
  m.choke = false
}

func (m *mockUpload) Choking() bool {
  return m.choke
}

func (m *mockUpload) IsInterested() bool {
  return m.isInterested
}

func (m *mockUpload) Rate() float64 {
  return m.rate
}

func (m *mockUpload) Handshake() Handshake {
  return nil
}

/* Tests. */

func makeUpload() (Upload, Transport, Transport, *Components) {
  t0 := NewTransferEngine(10, 10, "0")
  t1 := NewTransferEngine(10, 10, "1")

  c := &Components{
    Storage   : NewStorage("", []PieceMeta{
      PieceMeta{Index : 0, Length : 10},
    }, func() int {return 0}),
    Transport : t0,
    Choker    : &mockChoker{},
  }
  u := NewUpload(NewConnector("0", "1", c))
  return u, t0, t1, c
}

func TestChokeSendsMessageToPeer(t *testing.T) {
  u, _, t1, _ := makeUpload()

  assertEqual(t, u.Choking(), true)

  u.Unchoke()
  assertEqual(t, u.Choking(), false)
  assertEqual(t, <-t1.ControlRecv(), Unchoke{"0"})

  u.Choke()
  assertEqual(t, u.Choking(), true)
  assertEqual(t, <-t1.ControlRecv(), Choke{"0"})
}

func TestInterestedLetsChokerKnow(t *testing.T) {
  u, _, _, c := makeUpload()

  assertEqual(t, c.Choker.(*mockChoker).interestedCalled, false)
  u.Recv(Interested{"0"})
  assertEqual(t, c.Choker.(*mockChoker).interestedCalled, true)

  assertEqual(t, c.Choker.(*mockChoker).notInterestedCalled, false)
  u.Recv(NotInterested{"0"})
  assertEqual(t, c.Choker.(*mockChoker).notInterestedCalled, true)
}

// func TestUploadOnRequest(t *testing.T)
// func TestChokeStopsUploads(t *testing.T)
