package main

import (
  "github.com/danalex97/nfsTorrent/simulation"
  "github.com/danalex97/nfsTorrent/config"
  "github.com/danalex97/nfsTorrent/log"

  "flag"
  "math/rand"
  "time"
  "sync"
  "fmt"
  "os"
)

// Flags
var confPath = flag.String(
  "conf",
  "./config/confs/small.json",
  "The path to configuration .json file.",
)

var extension = flag.Bool(
  "ext",
  false,
  "Whether we use the extension",
)

func main() {
  // Parsing the flags
  flag.Parse()

  // Random seed
  rand.Seed(time.Now().UTC().UnixNano())

  var wg sync.WaitGroup

  var template interface {}
  if !*extension {
    template = new(simulation.SimulatedNode)
  } else {
    template = new(simulation.SimulatedCachedNode)
    fmt.Println("Running with extension.")
  }

  s := simulation.NewSimulation(
    template,
    config.
      JSONConfig(*confPath).
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
