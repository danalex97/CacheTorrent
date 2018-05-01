package main

import (
  "github.com/danalex97/nfsTorrent/simulation"
  "github.com/danalex97/nfsTorrent/config"
  "github.com/danalex97/nfsTorrent/log"

  errLog "log"
  "runtime/pprof"
  "runtime"

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
  "./confs/small.json",
  "The path to configuration .json file.",
)

var extension = flag.Bool(
  "ext",
  false,
  "Whether we use the extension",
)

var biased = flag.Bool(
  "bias",
  false,
  "Whether we use the biased tracker",
)

var verbose = flag.Bool(
  "v",
  false,
  "Verbose output",
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func cpuprofileRun() {
  if *cpuprofile != "" {
    f, err := os.Create(*cpuprofile)
    if err != nil {
        errLog.Fatal("could not create CPU profile: ", err)
    }
    if err := pprof.StartCPUProfile(f); err != nil {
        errLog.Fatal("could not start CPU profile: ", err)
    }
    defer pprof.StopCPUProfile()
  }
}

func memprofileRun() {
  if *memprofile != "" {
    f, err := os.Create(*memprofile)
    if err != nil {
        errLog.Fatal("could not create memory profile: ", err)
    }
    runtime.GC() // get up-to-date statistics
    if err := pprof.WriteHeapProfile(f); err != nil {
        errLog.Fatal("could not write memory profile: ", err)
    }
    f.Close()
  }
}

func main() {
  // Parsing the flags
  flag.Parse()

  // Random seed
  rand.Seed(time.Now().UTC().UnixNano())

  // Set verbosity
  log.SetVerbose(*verbose)

  // Profiling
  cpuprofileRun()

  var wg sync.WaitGroup

  // Set extension
  var template interface {}
  if !*extension {
    if !*biased {
      template = new(simulation.SimulatedNode)
    } else {
      template = new(simulation.SimulatedBiasedNode)
    }
  } else {
    template = new(simulation.SimulatedCachedNode)
    fmt.Println("Running with extension.")
  }

  // Run with configuration
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

  log.Query(log.GetRedundancy)
  log.Query(log.GetTraffic)
  log.Query(log.GetTime)
  log.Query(log.Stop)

  // Profiling
  memprofileRun()

  os.Exit(0)
}
