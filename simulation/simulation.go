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
    WithProgress(
      config.Config.AllNodesRun,
      config.Config.AllNodesRunInterval).
    //====================================
    WithCapacities().
    WithTransferInterval(
      config.Config.TransferInterval)

  if config.Config.Latency {
    builder = builder.WithLatency()
  }

  for _, tuple := range config.Config.CapacityNodes {
    builder = builder.WithCapacityNodes(
      tuple.Number,
      tuple.Upload,
      tuple.Download)
  }

  return builder.Build()
}
