package simulation

import (
  . "github.com/danalex97/Speer/sdk/go"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/model"
)

// the interface needed to be implemented by a peer
type Peer interface {
  New(util PeerUtil) Peer
  // the peer constructor

  OnJoin()
  // a method that should be called when a node joins the network

  OnLeave()
  // a method that should be called when a node leaves the network
}

// the interface that can be used by a peer
type PeerUtil interface {
  Id()        string
  // returns the id of the peer

  Transport() Transport
  // returns the transport interface of a peer
}

type SimulatedPeer interface {
  TorrentNode
  PeerUtil
}

type simulatedPeer struct {
  AutowiredTorrentNode

  id         string
  transport  Transport
  peer       Peer
}

/* Methods implemented for the TorrentNode interface. */
func (s *simulatedPeer) UnreliableNode() overlay.UnreliableNode {
  return nil
}

func (s *simulatedPeer) NewDHTNode() DHTNode {
  return nil
}

func (s *simulatedPeer) OnQuery(query model.DHTQuery) error {
  return nil
}

func (s *simulatedPeer) Key() string {
  return RandomKey()
}

/* Methods implemeneted for the PeerUtil interface. */
func (s *simulatedPeer) Id() string {
  return s.id
}

func (s *simulatedPeer) Transport() Transport {
  return s.transport
}

/* Methods implemeneted for the PeerUtil interface. */
func (s *simulatedPeer) OnJoin() {
  s.peer.OnJoin()
}

func (s *simulatedPeer) OnLeave() {
  s.peer.OnLeave()
}

func newSimulatedPeer(id string, transport Transport, peer Peer) SimulatedPeer {
  s := new(simulatedPeer)

  s.id        = id
  s.transport = transport

  // this will call the Peer constructor
  s.peer = peer.New(s)

  return s
}
