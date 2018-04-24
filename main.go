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

  s := simulation.NewSimulation(
    // new(simulation.SimulatedNode),
    new(simulation.SimulatedCachedNode),
    simulation.
      SmallTorrentConfig().
      WithParams(func(c *config.Conf) {
        c.SharedInit = func() {
          wg.Add(1)
        }
        c.SharedCallback = func() {
          wg.Done()
        }
      }),
  )

  // s := simulation.NewITLSimulation(
  //   // new(simulation.SimulatedNode),
  //   new(simulation.SimulatedCachedNode),
  //   simulation.
  //     ITLConfig().
  //     WithParams(func(c *config.Conf) {
  //       c.SharedInit = func() {
  //         wg.Add(1)
  //       }
  //       c.SharedCallback = func() {
  //         wg.Done()
  //         fmt.Println(wg)
  //       }
  //     }),
  // )

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
