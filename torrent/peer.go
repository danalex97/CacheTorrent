package torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "runtime"
  "fmt"
)

type Peer struct {
  id      string
  tracker string
  ids     []string

  transport Transport

  // used only to identify tracker
  join    string
}

func (p *Peer) OnJoin() {
  // If the transport is missing, it must be we are
  // on torrent-less node
  if p.transport == nil {
    return
  }

  // Find out who the tracker is
  p.transport.ControlSend(p.join, trackerReq{p.id})

  msg := <-p.transport.ControlRecv()
  p.tracker = msg.(trackerRes).id

  // The peer should be initialized
  fmt.Printf("Node %s started with tracker %s\n", p.id, p.tracker)

  // Send join message to the tracker
  p.transport.ControlSend(p.tracker, join{p.id})

  go func() {
    for {
      select {
      case m, ok := <-p.transport.ControlRecv():
        if !ok {
          continue
        }

        switch msg := m.(type) {
        case trackerReq:
          p.transport.ControlSend(msg.from, trackerRes{p.tracker})
        case neighbours:
          p.ids = msg.ids
          fmt.Println("Joined torrent: ", p.ids)
        }

      default:
        // allow other nodes in simulation run
        runtime.Gosched()
      }
    }
  }()
}

func (p *Peer) OnLeave() {
}

func (p *Peer) New(util TorrentNodeUtil) TorrentNode {
  peer := new(Peer)

  peer.id        = util.Id()
  peer.join      = util.Join()
  peer.transport = util.Transport()

  return peer
}
