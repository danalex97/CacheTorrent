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
  conns     []*Connector
}

func NewChoker(time func() int) *Choker {
  return &Choker{
    new(sync.Mutex),
    time,
    []*Connector{},
  }
}

func (c *Choker) AddConnector(conn *Connector) {
  c.Lock()
  defer c.Unlock()

  c.conns = append(c.conns, conn)

  // Send haves at connection
  s := conn.components.Storage
  t := conn.components.Transport
  for _, piece := range s.pieces {
    t.ControlSend(conn.to, have{conn.from, piece.index})
  }
}

func (c *Choker) rechoke() {
  c.Lock()
  defer c.Unlock()

  // We only upload to interested peers
  interested := []*Connector{}
  for _, conn := range c.conns {
    if conn.upload.isInterested {
      interested = append(interested, conn)
    }
  }

  // Sort the choked connections
  sort.Slice(interested, func(i, j int) bool {
    return interested[i].Rate() > interested[j].Rate()
  })

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

func (c *Choker) Interested(conn *Connector) {
  if !conn.upload.choke {
    c.rechoke()
  }
}

func (c *Choker) NotInterested(conn *Connector) {
  if !conn.upload.choke {
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

/**
 * We moved the responsibility of 'MultiDownload.py' to 'download.py'
 * and the functions below in the Choker as we only need a struct
 * which references the list of connections.
 */
func (c *Choker) Lost() {
  c.Lock()
  defer c.Unlock()

  for _, conn := range c.conns {
    // We try to request more pieces only if the connection is not choked
    if !conn.choked {
      conn.RequestMore()
    }
  }
}

func (c *Choker) Have(index int) {
  c.Lock()
  defer c.Unlock()

  for _, conn := range c.conns {
    t := conn.components.Transport
    t.ControlSend(conn.to, have{conn.from, index})
  }
}
