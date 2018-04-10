package main

import (
  "github.com/danalex97/nfsTorrent/simulation"
  "math/rand"
  "time"
  "fmt"
  "os"
)

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  s := simulation.SmallTorrentSimulation()
  s.Run()

  time.Sleep(time.Second * 3)
  fmt.Println("Done")
  s.Stop()

  os.Exit(0)
}
