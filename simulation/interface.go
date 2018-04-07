package simulation

import (
  . "github.com/danalex97/Speer/sdk/go"
  . "github.com/danalex97/Speer/model"
  "sync"
)

var mutex = new(sync.Mutex)
var ctr   = 0

var peer    Peer
var tracker Tracker

type SimulationTorrent struct {
  AutowiredTorrentNode

  id         string
  isTracker  bool

  node       TorrentNode
}

/* Interface functions. */
func (t *SimulationTorrent) OnJoin() {
  /* Wait for engine access. */
  Wait(func () bool {
    return t.Transfer() == nil
  })

  mutex.Lock()
  if ctr == 0 {
    t.isTracker = true
  } else {
    t.isTracker = false
  }
  ctr++
  mutex.Unlock()

  if t.isTracker {
    t.node = peer.New(t.id)
  } else {
    t.node = tracker.New(t.id)
  }

  t.node.OnJoin()
}

func (t *SimulationTorrent) OnQuery(query DHTQuery) error {
  if t.node == nil {
    return nil
  }
  return t.node.OnQuery(query)
}

func (t *SimulationTorrent) OnLeave() {
  if t.node != nil {
    t.node.OnLeave()
  }
}

func (t *SimulationTorrent) NewDHTNode() DHTNode {
  node := new(SimulationTorrent)

  node.Autowire(t)

  /* Initialize node. */
  node.id = t.UnreliableNode().Id()
  /* End initialize node. */

  return node
}

func (t *SimulationTorrent) Key() string {
  if t.node == nil {
    return ""
  }
  return t.node.Key()
}
