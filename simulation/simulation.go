package simulation

import (
  "github.com/danalex97/Speer/sdk/go"
  "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/config"
)

type Simulation interfaces.ISimulation

func SmallTorrentConfig() *config.Conf {
  return &config.Conf{
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

    SharedInit     : func() {},
    SharedCallback : func() {},
  }
}

func NewSimulation(template interface {}, newConfig *config.Conf) Simulation {
  if newConfig != nil {
    config.Config = newConfig
  }

  return sdk.NewDHTSimulationBuilder(template).
    WithPoissonProcessModel(2, 2).
    // WithInternetworkUnderlay(10, 50, 20, 50).
    WithInternetworkUnderlay(10, 50, 2, 50).
    WithDefaultQueryGenerator().
    WithLimitedNodes(config.Config.MinNodes + 1).
    // WithMetrics().
    //====================================
    WithCapacities().
    WithTransferInterval(10).
    WithCapacityNodes(config.Config.MinNodes + 1, 10, 20).
    Build()
}

/**
 * We want to do a similar simulation to the ones done in:
 * [R. Bindal et al., "Improving Traffic Locality in BitTorrent via
 * Biased Neighbor Selection," 26th IEEE International Conference on
 * Distributed Computing Systems (ICDCS'06), 2006, pp. 66-66]
 *
 * - topology:
 *   - 700 peers
 *   - 14 ISPs
 *   - around 50 peers/ISP
 * - links:
 *   - unit of measure: b/ms
 *   - uplink:   400Kb/s = 0.4Kb/ms = 400 b/ms
 *   - downlink: 1.5Mb/s = 1,500Kb/s = 1.5Kb/ms = 1,500 b/ms
 * - BitTorrent configuration:
 *   - out peers: 35
 *   - rechoking interval: 10s = 10,000ms
 *   - 5 unchoked connections with 1 optimistic
 *   - piece size: 245KB = 245,000B = 1,960,000b
 *
 *   - pieces: 1000 (default)
 *   - backlog: 10 (default)
 */
func ITLConfig() *config.Conf {
 return &config.Conf{
   OutPeers : 35,
   InPeers  : 35,

   MinNodes : 700,
   Seeds    : 1,

   PieceSize : 1960000,
   Pieces    : 1000,

   Uploads     : 4,
   Optimistics : 1,
   Interval    : 10000,

   Backlog : 10,

   SharedInit     : func() {},
   SharedCallback : func() {},
 }
}

func NewITLSimulation(template interface {}, newConfig *config.Conf) Simulation {
  if newConfig != nil {
    config.Config = newConfig
  }

  return sdk.NewDHTSimulationBuilder(template).
    WithPoissonProcessModel(2, 2).
    // transitDomains, transitDomainSize, stubDomains, stubDomainSize
    WithInternetworkUnderlay(10, 50, 14, 100).
    WithDefaultQueryGenerator().
    WithLimitedNodes(config.Config.MinNodes + 1).
    //====================================
    WithCapacities().
    // unit: ms
    WithTransferInterval(100).
    // number, up, down; unit: b/ms
    WithCapacityNodes(config.Config.MinNodes + 1, 400, 1500).
    Build()
}
