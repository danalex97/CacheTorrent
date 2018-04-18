package cache_torrent

import (
  "strings"
)

type Election struct {
}

func NewElection() *Election {
  return &Election{}
}

func (e *Election) Run() {
}

func (e *Election) Recv(m interface {}) {
}

func (e *Election) NewJoin(id string) {

}

func getAS(id string) string {
  // We assume that ID is of form [AS].[NBR]
  return strings.Split(id, ".")[0]
}
