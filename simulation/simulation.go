package simulation

import (
  "github.com/danalex97/Speer/sdk/go"
  "github.com/danalex97/Speer/interfaces"

  "github.com/danalex97/nfsTorrent/config"
)

type Simulation interfaces.ISimulation

func SmallTorrentSimulation(template interface {}, newConfig *config.Conf) Simulation {
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
