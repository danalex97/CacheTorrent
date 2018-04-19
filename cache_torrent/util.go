package cache_torrent

import (
  "strings"
)

func getNBR(id string) string {
  // We assume that ID is of form [AS].[NBR]
  return strings.Split(id, ".")[0]
}

func getAS(id string) string {
  // We assume that ID is of form [AS].[NBR]
  return strings.Split(id, ".")[0]
}

func isLeader(p *Peer) bool {
  // Check if peer is leader
  for _, id := range p.Leaders {
    if id == p.Id {
      return true
    }
  }
  return false
}
