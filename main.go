package main

import (
  "github.com/danalex97/nfsTorrent/simulation"
  "github.com/danalex97/nfsTorrent/config"
  "github.com/danalex97/nfsTorrent/log"

  errLog "log"
  "runtime/pprof"
  "runtime"
  "os/signal"
  "os"

  "flag"
  "math/rand"
  "time"
  "sync"
  "fmt"
)

const MaxUint = ^uint(0)
const MaxInt  = int(MaxUint >> 1)

// Flags
var confPath = flag.String(
  "conf",
  "./confs/small.json",
  "The path to configuration .json file.",
)

var extension = flag.Int(
  "ext",
  MaxInt,
  "Use the textesion with ext percent number of leaders.",
)

var timeCDF = flag.Bool(
  "cdf",
  false,
  "Enable printing time cumulative distribution function.",
)

var biased = flag.Int(
  "bias",
  MaxInt,
  "Number of outgoing connections for a biased Tracker.",
)

var verbose = flag.Bool(
  "v",
  false,
  "Verbose output",
)

var leaders = flag.Bool(
  "leaders",
  false,
  "Enable printing leader and follower times.",
)

var pieces = flag.Int(
  "pieces",
  MaxInt,
  "Number of pieces the file has.",
)

var pieceSize = flag.Int(
  "pieceSize",
  MaxInt,
  "The size of a piece from the file.",
)

var logfile = flag.String(
  "log",
  "",
  "The packet log in `.json` format.",
)

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to `file`.")
var memprofile = flag.String("memprofile", "", "Write memory profile to `file`.")

func makeMemprofile() {
  // Profiling
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

func getStats() {
  log.Query(log.GetRedundancy)
  log.Query(log.GetTraffic)
  log.Query(log.GetTime)

  if *extension != MaxInt &&  *leaders {
    log.Query(log.GetTimeLeader)
  }

  if *timeCDF {
    log.Query(log.GetTimeCDF)
    if *extension != MaxInt &&  *leaders {
      log.Query(log.GetLeaderCDF)
    }
  }

  log.Query(log.Stop)
}

func main() {
  fmt.Println("Running:", os.Args[1:])

  // Parsing the flags
  flag.Parse()

  // Profiling
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
  defer makeMemprofile()
  // Get profile even on signal
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  go func(){
    for sig := range c {
      fmt.Println("Singal received:", sig)
      if *cpuprofile != "" {
        pprof.StopCPUProfile()
      }
      makeMemprofile()

      fmt.Println("Partial stats:")
      getStats()
      os.Exit(0)
    }
  }()

  // Random seed
  rand.Seed(time.Now().UTC().UnixNano())

  // Set verbosity
  log.SetVerbose(*verbose)

  // Set log file
  log.SetLogfile(*logfile)

  var wg sync.WaitGroup

  // Set extension
  var template interface {}
  if *extension == MaxInt {
    if *biased == MaxInt {
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
        c.Bias          = *biased
        c.LeaderPercent = *extension

        if *pieces != MaxInt {
          c.Pieces = *pieces
        }

        if *pieceSize != MaxInt {
          c.PieceSize = *pieceSize
        }

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

  getStats()
}
