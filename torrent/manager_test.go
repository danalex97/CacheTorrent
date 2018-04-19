package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  . "github.com/danalex97/Speer/capacity"
  "testing"
)

/* Mocks. */
type mockManager struct {
  uploads   []Upload
  downloads []Download
}

func (m *mockManager) AddConnector(conn *Connector) {}
func (m *mockManager) Uploads()   []Upload    { return m.uploads }
func (m *mockManager) Downloads() []Download { return m.downloads }

/* Tests. */
func TestManagerConcurrent(t *testing.T) {
  for i := 0; i < 10; i++ {
    m := NewConnectionManager()
    s := NewStorage("", []PieceMeta{})

    s.Store(Piece{"", 0, 0, Data{"0", 10}})

    conns := []*Connector{
      NewConnector("0", "1", &Components{
        Storage   : s,
        Transport : NewTransferEngine(10, 10, "0"),
      }).WithUpload().WithDownload().WithHandshake(),
      NewConnector("1", "0", &Components{
        Storage   : s,
        Transport : NewTransferEngine(10, 10, "1"),
      }).WithUpload().WithDownload().WithHandshake(),
    }

    done := make(chan bool)
    add  := func(conn *Connector) {
      m.AddConnector(conn)
      done <- true
    }
    go add(conns[0])
    go add(conns[1])

    for j := 0; j < 2; j++ {
      <- done
    }

    // Check have messages gets sent
    assertEqual(t, <-conns[0].Transport.ControlRecv(), Have{"1", 0})
    assertEqual(t, <-conns[1].Transport.ControlRecv(), Have{"0", 0})

    // Check the number of connections is 2
    go assertEqual(t, len(m.Uploads()), 2)
    go assertEqual(t, len(m.Uploads()), 2)
    go assertEqual(t, len(m.Downloads()), 2)
    go assertEqual(t, len(m.Downloads()), 2)
  }
}
