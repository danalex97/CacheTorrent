package simulation

import (
  . "github.com/danalex97/Speer/sdk/go"
  . "github.com/danalex97/Speer/model"
  "sync"
)

var mutex = new(sync.Mutex)
var ctr   = 0

var peer    Peer
var tracker Peer

type SimulatedTorrent struct {
  AutowiredTorrentNode

  id         string
  isTracker  bool

  node       TorrentNode
}

/* Interface functions. */
func (t *SimulatedTorrent) OnJoin() {
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
    t.node = newSimulatedPeer(t.id, t.Transfer(), peer)
  } else {
    t.node = newSimulatedPeer(t.id, t.Transfer(), tracker)
  }

  t.node.OnJoin()
}

func (t *SimulatedTorrent) OnQuery(query DHTQuery) error {
  if t.node == nil {
    return nil
  }
  return t.node.OnQuery(query)
}

func (t *SimulatedTorrent) OnLeave() {
  if t.node != nil {
    t.node.OnLeave()
  }
}

func (t *SimulatedTorrent) NewDHTNode() DHTNode {
  node := new(SimulatedTorrent)

  node.Autowire(t)

  /* Initialize node. */
  node.id = t.UnreliableNode().Id()
  /* End initialize node. */

  return node
}

func (t *SimulatedTorrent) Key() string {
  if t.node == nil {
    return ""
  }
  return t.node.Key()
}
