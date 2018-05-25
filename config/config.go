package config

type NodeConf struct {
  Number    int `json:"number"`
  Upload    int `json:"upload"`
  Download  int `json:"download"`
}

type Conf struct {
  // BitTorrent
  OutPeers int `json:"outPeers"` // maximum number of outbound peers
  InPeers  int `json:"inPeers"`  // maximum number of inbound peers

  MinNodes int `json:"minNodes"` // minimum number of nodes in a Torrent
  Seeds    int `json:"seeds"`    // number of seed nodes

  PieceSize int `json:"pieceSize"` // the piece size
  Pieces    int `json:"pieces"`    // number of pieces in a Torrent

  Uploads     int `json:"uploads"`    // number of Uploads (without Optimistics)
                                      //   chosen by a Choker
  Optimistics int `json:"optimistics"`// number of Optimistics
                                      //   chosen by a Choker
  Interval    int `json:"interval"`   // milliseconds

  Backlog int `json:"backlog"`// number of pieces requested at a time

  // Extension
  LeaderPercent int `json:"leaderPercent"`// the percent of leaders in an AS

  // Biased Tracker extension
  Bias int `json:"bias"` // the number of connections to peers in different AS

  // Misc
  SharedCallback func()
  SharedInit     func()

  // Simulation parameters
  TransitDomains     int `json:"transitDomains"`
  TransitDomainSize  int `json:"transitDomainSize"`
  StubDomains        int `json:"stubDomains"`
  StubDomainSize     int `json:"stubDomainSize"`

  TransferInterval   int `json:"transferInterval"`

  CapacityNodes      []NodeConf `json:"capacityNodes"`

  // Progress properties
  AllNodesRun         *WGProgress
  AllNodesRunInterval int `json:allNodesRunInterval`

  // Latency support
  Latency bool

  // Run the simulator's event queue with support for parallel events.
  Parallel bool
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
