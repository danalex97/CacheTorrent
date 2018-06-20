package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/nfsTorrent/config"

  "testing"
)

/* Mocks. */
type mockTorrentNodeUtil struct {
  t    Transport
  id   string
  join string
}

func (u *mockTorrentNodeUtil) Id() string {
  return u.id
}

func (u *mockTorrentNodeUtil) Join() string {
  return u.join
}

func (u *mockTorrentNodeUtil) Time() func() int {
  return func() int {
    return 0
  }
}

func (u *mockTorrentNodeUtil) Transport() Transport {
  return u.t
}

type mockControlTransport struct {
  recv chan interface{}
  send map[string]chan interface{}

  up   int
  down int
}

func newMockControlTransport() *mockControlTransport {
  return &mockControlTransport{
    recv : make(chan interface{}, 10),
    send : make(map[string]chan interface{}),

    up   : 10,
    down : 10,
  }
}

func (t *mockControlTransport) init(id string) {
  t.send[id] = make(chan interface{}, 10)
}

func (t *mockControlTransport) Up() int {
  return t.up
}

func (t *mockControlTransport) Down() int {
  return t.down
}

func (t *mockControlTransport) Connect(_ string) Link {
  return nil
}

func (t *mockControlTransport) ControlPing(_ string) bool {
  return true
}

func (t *mockControlTransport) ControlSend(id string, message interface {}) {
  t.send[id] <- message
}

func (t *mockControlTransport) ControlRecv() <-chan interface{} {
  return t.recv
}

func newPeer(id, join string) (*Peer, *mockControlTransport) {
  t := newMockControlTransport()
  t.init(id)
  t.init(join)
  return new(Peer).New(&mockTorrentNodeUtil{
    id   : id,
    join : join,

    t : t,
  }).(*Peer), t
}

/* Tests. */
func TestInitRespondsToAllTrackerReq(t *testing.T) {
  p, tr := newPeer("1", "0")

  go p.Init()

  // Check peer sends tracker request at joining
  assertEqual(t, <-tr.send["0"], TrackerReq{"1"})

  // Send TrackerReq
  tr.init("2")
  tr.init("3")
  tr.recv <- TrackerReq{"2"}
  tr.recv <- TrackerReq{"3"}

  // Response from Tracker
  tr.recv <- TrackerRes{"0"}

  // Check responses get sent
  assertEqual(t, <-tr.send["2"], TrackerRes{"0"})
  assertEqual(t, <-tr.send["3"], TrackerRes{"0"})

  // Check join was sent
  assertEqual(t, <-tr.send["0"], Join{"1"})
}

func TestBindReturnsCorrectState(t *testing.T) {
  p, tr := newPeer("1", "0")
  p.Tracker = "0"

  tr.init("2")
  assertEqual(t, p.Bind(TrackerReq{"2"}), BindDone)
  assertEqual(t, p.Bind(Neighbours{[]string{"2"}}), BindDone)

  assertEqual(t, p.Bind(SeedRes{[]PieceMeta{}}), BindRun)

  assertEqual(t, p.Bind(nil), BindNone)

  p.Connectors["1"] = &mockRunner{}
  assertEqual(t, p.Bind(nil), BindRecv)
}

func TestRunInitAllComponentsAndConnections(t *testing.T) {
  p, _ := newPeer("1", "0")

  config.Config = new(config.Conf)

  total := 0
  add := func(_ string) {
    total += 1
  }

  p.Ids = []string{"2", "3", "4"}
  p.Run(add)

  assertEqual(t, total, 3)

  if p.Storage == nil {
    t.Fatalf("Storage not initialized.")
  }
  if p.Picker == nil {
    t.Fatalf("Picker not initialized.")
  }
  if p.Manager == nil {
    t.Fatalf("Manager not initialized.")
  }
  if p.Choker == nil {
    t.Fatalf("Choker not initialized.")
  }
}

func TestRunRecvAddsConnection(t *testing.T) {
  p, _ := newPeer("1", "0")

  total := 0
  add := func(_ string) {
    total += 1
  }

  p.RunRecv("1", nil, add)
  assertEqual(t, total, 1)
}
