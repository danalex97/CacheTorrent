package torrent

import (
  "github.com/danalex97/nfsTorrent/config"
  "math/rand"
  "runtime"
  "sort"
  "sync"
)

// We use variables instead of constants to allow testing.
var uploads     config.Const = config.NewConst(config.Uploads)
var optimistics config.Const = config.NewConst(config.Optimistics)
var interval    config.Const = config.NewConst(config.Interval)

var strategy config.Const = config.NewConst(config.ChokerStrategy)

// A Choker is a structure which periodically chokes and unchokes connections
// based on the Tit-for-Tat strategy. When an Upload component is unchoked,
// the respective upload component is allowed to upload new pieces to other
// peers. Furthermore, the Choker reacts to Interested and NotInterested
// notifications arriving from Upload components.
//
// We consider only the rate heuristic used for peers rather than
// the heuristic used for seeds. That is, we use the connections with
// the highest upload rate.
//
// We do not model snubbing.
//
// For modelling simplicity, we randomize the optimistics.
//
// Version 5.3.0 uses allocates 30% of the slots to seeds and 70% to other
// peers.
//
// The seeder will have no upload made to it. The seeder can either:
//  - upload to best download rates
//  - upload randomly
// For the moment our implementation is random.
type Choker interface {
  Interested(conn Upload)
  NotInterested(conn Upload)

  Run()
}

type choker struct {
  *sync.Mutex

  time      func() int
  manager   Manager
}

func NewChoker(manager Manager, time func() int) Choker {
  return &choker{
    Mutex:    new(sync.Mutex),

    time:     time,
    manager:  manager,
  }
}

type byRate []Upload

func (a byRate) Len() int           { return len(a) }
func (a byRate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byRate) Less(i, j int) bool { return a[i].Rate() > a[j].Rate() }

func (c *choker) rechoke() {
  c.Lock()
  defer c.Unlock()

  conns := c.manager.Uploads()

  // We only upload to interested peers
  interested := []Upload{}
  for _, conn := range conns {
    if conn.IsInterested() {
      interested = append(interested, conn)
    }
  }

  if strategy.String() == config.Random {
    // -- Random strategy --
    // We choose the optimistics randomly for simplicity of modelling
    unchoked := uploads.Int() + optimistics.Int()
    if unchoked > len(interested) {
      unchoked = len(interested)
    }

    // Unchoke randomly from the interested connections
    perm := rand.Perm(len(interested))
    for i := 0; i < unchoked; i++ {
      interested[perm[i]].Unchoke()
    }
    for i := unchoked; i < len(interested); i++ {
      interested[perm[i]].Choke()
    }
  } else {
    // -- TitForTat strategy --

    // Sort the choked connections
    sort.Sort(byRate(interested))

    // If we want to consider the seeds, we should use 2 separate lists.

    // Unchoke the pereferred connections
    unchoked := uploads.Int()
    if unchoked > len(interested) {
      unchoked = len(interested)
    }
    for i := 0; i < unchoked; i++ {
      interested[i].Unchoke()
    }

    // Chocke the rest and handle optimistics
    rest := interested[unchoked:]
    unchoked = optimistics.Int()
    if unchoked > len(rest) {
      unchoked = len(rest)
    }
    // We choose the optimistics randomly for simplicity of modelling
    perm := rand.Perm(len(rest))
    for i := 0; i < unchoked; i++ {
      rest[perm[i]].Unchoke()
    }
    for i := unchoked; i < len(rest); i++ {
      rest[perm[i]].Choke()
    }
  }
}

func (c *choker) Interested(conn Upload) {
  if !conn.Choking() {
    c.rechoke()
  }
}

func (c *choker) NotInterested(conn Upload) {
  if !conn.Choking() {
    c.rechoke()
  }
}

func (c *choker) Run() {
  c.rechoke()
  t := c.time()
  l := t
  for {
    t = c.time()

    // This seems to work fine for up to 1000 nodes.
    if t - l > interval.Int() {
      c.rechoke()
      l = t
    }
    runtime.Gosched()
  }
}
