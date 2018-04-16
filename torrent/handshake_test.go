package torrent

import (
  . "github.com/danalex97/Speer/capacity"
  "testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
    t.Fatalf("%s != %s", a, b)
	}
}

func TestHandshakeBidir(t *testing.T) {
  for i := 0; i < 10; i++ {
    conn0 := NewConnector("0", "1", &Components{
      Transport : NewTransferEngine(10, 10, "0"),
    })
    conn1 := NewConnector("1", "0", &Components{
      Transport : NewTransferEngine(10, 10, "1"),
    })

    h0 := NewHandshake(conn0).(*handshake)
    h1 := NewHandshake(conn1).(*handshake)

    go h0.Run()
    go h1.Run()

    go h0.Recv(<-h0.Transport.ControlRecv())
    go h1.Recv(<-h1.Transport.ControlRecv())

    assertEqual(t, h1.Uplink(), h0.Downlink())
    assertEqual(t, h0.Uplink(), h1.Downlink())
    assertEqual(t, h1.Done(), true)
    assertEqual(t, h0.Done(), true)
  }
}

func TestHandshakeUnidir(t *testing.T) {
  for i := 0; i < 10; i++ {
    conn0 := NewConnector("0", "1", &Components{
      Transport : NewTransferEngine(10, 10, "0"),
    })
    conn1 := NewConnector("1", "0", &Components{
      Transport : NewTransferEngine(10, 10, "1"),
    })

    h0 := NewHandshake(conn0).(*handshake)
    h1 := NewHandshake(conn1).(*handshake)

    // start h0 handshake
    go h0.Run()

    // wait using h1 handshake
    go h1.Recv(<-h1.Transport.ControlRecv())

    assertEqual(t, h0.Uplink(), h1.Downlink())
    assertEqual(t, h1.Done(), true)
    assertEqual(t, h0.Done(), false)
  }
}
