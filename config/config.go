package config

type Conf struct {
  OutPeers int // maximum number of outbound peers
  InPeers  int // maximum number of inbound peers

  MinNodes int // minimum number of nodes in a Torrent
  Seeds    int // number of seed nodes

  PieceSize int // the piece size
  Pieces    int // number of pieces in a Torrent

  Uploads     int // number of Uploads (without Optimistics) chosen by a Choker
  Optimistics int // number of Optimistics chosen by a Choker
  Interval    int // milliseconds

  Backlog        int // number of pieces requested at a time

  Shared         interface {}
  SharedCallback func(interface {})
}

func NewConf() *Conf {
  return &Conf{
    OutPeers : 3,
    InPeers  : 3,

    MinNodes : 10,
    Seeds    : 1,

    PieceSize : 10,
    Pieces    : 1,

    Uploads     : 0,
    Optimistics : 1,
    Interval    : 10000,

    Backlog : 10,

    Shared         : nil,
    SharedCallback : func(interface {}) {},
  }
}

/**
 * Usage:
 *  c := NewConf().
 *    WithParams(func(c *Conf) {
 *      c.OutPeers = 5
 *      c.MinNodes = 10
 *    })
 */
func (c *Conf) WithParams(f func (c *Conf)) *Conf {
  f(c)
  return c
}

var Config *Conf = NewConf()
