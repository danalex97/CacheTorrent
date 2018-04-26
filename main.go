package main

import (
  "github.com/danalex97/nfsTorrent/simulation"
  "github.com/danalex97/nfsTorrent/config"
  "github.com/danalex97/nfsTorrent/log"

  "math/rand"
  "time"
  "sync"
  "fmt"
  "os"
)

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  var wg sync.WaitGroup

  fmt.Println(config.JSONConfig("./config/confs/small.json"))

  s := simulation.NewSimulation(
    // new(simulation.SimulatedNode),
    new(simulation.SimulatedCachedNode),
    config.
      JSONConfig("./config/confs/small.json").
      // JSONConfig("./config/confs/itl.json").
      WithParams(func(c *config.Conf) {
        c.SharedInit = func() {
          wg.Add(1)
        }
        c.SharedCallback = func() {
          wg.Done()
        }
      }),
  )

  s.Run()

  // Initial time required to run SharedInits
  time.Sleep(time.Duration(float64(time.Second) * 1))
  fmt.Println("Init period done.")

  // Wait for all nodes to finish.
  wg.Done()
  wg.Wait()

  s.Stop()
  t := s.Time()
  fmt.Println("Downloads finished in", t, "milliseconds.")

  log.Log.Query(log.GetRedundancy)
  log.Log.Query(log.GetTime)
  log.Log.Query(log.Stop)

  os.Exit(0)
}
