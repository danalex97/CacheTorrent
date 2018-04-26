package simulation

import (
  "github.com/danalex97/Speer/sdk/go"
  "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/config"
)

type Simulation interfaces.ISimulation

func NewSimulation(template interface {}, newConfig *config.Conf) Simulation {
  config.Config = newConfig

  builder := sdk.NewDHTSimulationBuilder(template).
    WithPoissonProcessModel(2, 2).
    WithInternetworkUnderlay(
      config.Config.TransitDomains,
      config.Config.TransitDomainSize,
      config.Config.StubDomains,
      config.Config.StubDomainSize).
    WithDefaultQueryGenerator().
    WithLimitedNodes(config.Config.MinNodes + 1).
    // WithMetrics().
    //====================================
    WithCapacities().
    WithTransferInterval(
      config.Config.TransferInterval)

  for _, tuple := range config.Config.CapacityNodes {
    builder = builder.WithCapacityNodes(
      tuple.Number,
      tuple.Upload,
      tuple.Download)
  }

  return builder.Build()
}

func SmallTorrentConfig() *config.Conf {
  return &config.Conf{
    OutPeers : 10,
    InPeers  : 10,

    MinNodes : 20,
    Seeds    : 1,

    PieceSize : 196,
    Pieces    : 1000,

    Uploads     : 4,
    Optimistics : 1,
    Interval    : 10000,

    Backlog : 10,

    LeaderPercent : 50,

    SharedInit     : func() {},
    SharedCallback : func() {},

    TransitDomains    : 10,
    TransitDomainSize : 50,
    StubDomains       : 2,
    StubDomainSize    : 50,

    TransferInterval  : 10,

    CapacityNodes : []config.NodeConf{config.NodeConf{
      Number   : 20 + 1,
      Upload   : 10,
      Download : 20,
    }},
  }
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

   PieceSize : 196,
   Pieces    : 1000,

   Uploads     : 4,
   Optimistics : 1,
   Interval    : 10000,

   Backlog : 10,

   LeaderPercent : 50,

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
