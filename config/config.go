package config

type NodeConf struct {
  Number    int
  Upload    int
  Download  int
}

type Conf struct {
  // BitTorrent
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

  // Extension
  LeaderPercent int // the percent of leaders in an AS

  // Misc
  SharedCallback func()
  SharedInit     func()

  // Simulation parameters
  TransitDomains     int
  TransitDomainSize  int
  StubDomains        int
  StubDomainSize     int

  TransferInterval   int
  CapacityNodes      []NodeConf
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

var Config *Conf
