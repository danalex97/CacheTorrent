package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  . "github.com/danalex97/Speer/capacity"
  "testing"
)

/* Tests. */

func makeDownload() (Download, Transport, Transport, *Components) {
  backlog     = &mockConst{10}
  pieceNumber = &mockConst{10}

  t0 := NewTransferEngine(10, 10, "0")
  t1 := NewTransferEngine(10, 10, "1")

  s := NewStorage("", []PieceMeta{})
  c := &Components{
    Storage   : s,
    Picker    : NewPicker(s),
    Manager   : &mockManager{},
    Transport : t0,
    Choker    : &mockChoker{},
  }
  d := NewDownload(NewConnector("0", "1", c))
  return d, t0, t1, c
}

func TestGotHaveSendsInterestedMessageIfChoked(t *testing.T) {
  d, _, t1, c := makeDownload()

  c.Picker.GotHave("1", 0)
  d.Recv(Choke{"1"})

  assertEqual(t, d.Choked(), true)
  assertEqual(t, <-t1.ControlRecv(), Interested{"0"})
}

func TestGotHaveLetsPickerKnow(t *testing.T) {
  d, _, _, c := makeDownload()

  _, ok := c.Picker.Next("1")
  assertEqual(t, ok, false)

  d.Recv(Have{"1", 0})

  _, ok = c.Picker.Next("1")
  assertEqual(t, ok, true)
}

func TestGotChokeRequestsDoesntResendInterestChange(t *testing.T) {
  d, _, t1, _ := makeDownload()

  d.Recv(Have{"1", 0})
  assertEqual(t, d.Interested(), true)
  assertEqual(t, <-t1.ControlRecv(), Interested{"0"})

  d.Recv(Unchoke{"1"})
  assertEqual(t, d.Choked(), false)
  assertEqual(t, <-t1.ControlRecv(), Request{"0", 0})

  d.Recv(Choke{"1"})
  assertEqual(t, d.Choked(), true)
  assertEqual(t, d.Interested(), true)

  assertEqual(t, len(t1.ControlRecv()), 0)
}

func TestGotChokeSendsInterestedMessages(t *testing.T) {
  d, _, t1, c := makeDownload()

  assertEqual(t, d.Interested(), false)
  c.Picker.GotHave("1", 0)
  d.Recv(Choke{"1"})

  assertEqual(t, <-t1.ControlRecv(), Interested{"0"})
}

// func TestGotChokeRequestsLostActives(t *testing.T)

func TestGotUnchokeRequestsMore(t *testing.T) {
  d, _, t1, _ := makeDownload()

  d.Recv(Have{"1", 0})
  assertEqual(t, d.Interested(), true)
  assertEqual(t, <-t1.ControlRecv(), Interested{"0"})

  d.Recv(Unchoke{"1"})
  assertEqual(t, d.Choked(), false)
  assertEqual(t, <-t1.ControlRecv(), Request{"0", 0})
}

func TestGotPieceSendsHaves(t *testing.T) {
  d, _, t1, c := makeDownload()

  c.Manager.(*mockManager).uploads = []Upload{&mockUpload{
    me : "0",
    to : "1",
  }}
  d.Recv(Piece{"1", 0, 0, Data{"0", 10}})

  assertEqual(t, <-t1.ControlRecv(), Have{"0", 0})
}

func TestGotPieceRequestsMoreIfNotChoked(t *testing.T) {
  d, _, t1, c := makeDownload()

  d.(*TorrentDownload).choked = false

  c.Picker.GotHave("1", 1)
  d.Recv(Piece{"1", 0, 0, Data{"0", 10}})

  assertEqual(t, <-t1.ControlRecv(), Interested{"0"})
  assertEqual(t, <-t1.ControlRecv(), Request{"0", 1})
}

func TestGotPieceStoresPiece(t *testing.T) {
  d, _, _, c := makeDownload()

  _, ok := c.Storage.Have(0)
  assertEqual(t, ok, false)

  d.Recv(Piece{"1", 0, 0, Data{"0", 10}})

  _, ok = c.Storage.Have(0)
  assertEqual(t, ok, true)
}

// func TestRun(t *testing.T) {}
