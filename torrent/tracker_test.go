package torrent

import (
  "github.com/danalex97/nfsTorrent/config"

  "testing"
)

func newTracker(id string) (*Tracker, *mockControlTransport) {
  t := newMockControlTransport()
  t.init(id)

  return new(Tracker).New(&mockTorrentNodeUtil{
    id   : id,
    join : "",

    t : t,
  }).(*Tracker), t
}

func TestTracker(t *testing.T) {
  config.Config = &config.Conf{
    MinNodes  : 5,
    Seeds     : 1,
    Pieces    : 10,
    PieceSize : 10,
    OutPeers  : 4,
  }

  p, tr := newTracker("0")

  tr.init("1")
  tr.init("2")
  tr.init("3")
  tr.init("4")
  tr.init("5")

  p.OnJoin()

  tr.recv <- TrackerReq{"1"}

  tr.recv <- Join{"2"}
  tr.recv <- Join{"3"}
  tr.recv <- Join{"4"}
  tr.recv <- Join{"5"}
  tr.recv <- Join{"1"}

  tr.recv <- SeedReq{"1"}
  tr.recv <- SeedReq{"2"}

  assertEqual(t, <-tr.send["1"], TrackerRes{"0"})

  assertEqual(t, len((<-tr.send["1"]).(Neighbours).Ids), 4)
  assertEqual(t, len((<-tr.send["2"]).(Neighbours).Ids), 4)
  assertEqual(t, len((<-tr.send["3"]).(Neighbours).Ids), 4)
  assertEqual(t, len((<-tr.send["4"]).(Neighbours).Ids), 4)
  assertEqual(t, len((<-tr.send["5"]).(Neighbours).Ids), 4)

  assertEqual(t, len((<-tr.send["1"]).(SeedRes).Pieces), 0)
  assertEqual(t, len((<-tr.send["2"]).(SeedRes).Pieces), 10)
}
