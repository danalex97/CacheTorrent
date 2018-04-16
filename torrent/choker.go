package torrent

/*
 * We consider only the rate heuristic used for peers rather than
 * the heuristic used for seeds. That is, we use the connections with
 * the highest upload rate.
 *
 * We do not model snubbing.
 *
 * For modelling simplicity, we randomize the optimistics.
 *
 * Version 5.3.0 uses allocates 30% of the slots to seeds and 70% to other
 * peers.
 *
 * The seeder will have no upload made to it. The seeder can either:
 *  - upload to best download rates
 *  - upload randomly
 * For the moment our implementation is random.
 *
 */

import (
  "github.com/danalex97/nfsTorrent/config"
  "math/rand"
  "runtime"
  "sort"
  "sync"
)

const uploads     int = config.Uploads
const optimistics int = config.Optimistics
const interval    int = config.Interval

type Choker struct {
  *sync.Mutex

  time      func() int
  manager   Manager
}

func NewChoker(manager Manager, time func() int) *Choker {
  return &Choker{
    Mutex:    new(sync.Mutex),

    time:     time,
    manager:  manager,
  }
}

type byRate []Upload

func (a byRate) Len() int           { return len(a) }
func (a byRate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byRate) Less(i, j int) bool { return a[i].Rate() > a[j].Rate() }

func (c *Choker) rechoke() {
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

  // Sort the choked connections
  sort.Sort(byRate(interested))

  // If we want to consider the seeds, we should use 2 separate lists.

  // Unchoke the pereferred connections
  unchoked := uploads
  if unchoked > len(interested) {
    unchoked = len(interested)
  }
  for i := 0; i < unchoked; i++ {
    interested[i].Unchoke()
  }

  // Chocke the rest and handle optimistics
  rest := interested[unchoked:]
  unchoked = optimistics
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

func (c *Choker) Interested(conn Upload) {
  if !conn.Choking() {
    c.rechoke()
  }
}

func (c *Choker) NotInterested(conn Upload) {
  if !conn.Choking() {
    c.rechoke()
  }
}

func (c *Choker) Run() {
  c.rechoke()
  t := c.time()
  l := t
  for {
    t = c.time()

    // This seems to work fine for up to 1000 nodes.
    if t - l > interval {
      c.rechoke()
      l = t
    }
    runtime.Gosched()
  }
}
