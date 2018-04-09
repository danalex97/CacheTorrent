package simulation

import (
  "github.com/danalex97/Speer/sdk/go"
  "github.com/danalex97/Speer/interfaces"
)

type Simulation interfaces.ISimulation

func SmallTorrentSimulation() Simulation {
  nodeTemplate := new(simulatedNode)
  return sdk.NewDHTSimulationBuilder(nodeTemplate).
    WithPoissonProcessModel(2, 2).
    WithInternetworkUnderlay(10, 50, 20, 50).
    WithDefaultQueryGenerator().
    WithLimitedNodes(11).
    // WithMetrics().
    //====================================
    WithCapacities().
    WithTransferInterval(10).
    WithCapacityNodes(11, 10, 20).
    Build()
}
