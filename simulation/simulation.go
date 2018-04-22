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
