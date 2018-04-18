package cache_torrent

import (
  "sync"
  "sort"
  "strings"
)

const MaxLeaders int = 1

type Election struct {
  sync.Mutex

  camera      map[string][]string
  candidates  map[string][]Candidate
}

func NewElection() *Election {
  return &Election{
    camera     : make(map[string][]string),
    candidates : make(map[string][]Candidate),
  }
}

func (e *Election) Run() {
}

func (e *Election) Recv(m interface {}) {
  switch candidate := m.(type) {
  case Candidate:
    e.RegisterCandidate(candidate)
  }
}

func (e *Election) RegisterCandidate(candidate Candidate) {
  e.Lock()
  defer e.Unlock()

  /* Add candidate to candidate list. */
  as := getAS(candidate.Id)
  if _, ok := e.candidates[as]; !ok {
    e.candidates[as] = []Candidate{}
  }
  e.candidates[as] = append(e.candidates[as], candidate)
}

func (e *Election) NewJoin(id string) {
  e.Lock()
  defer e.Unlock()

  /* Add id to camera. */
  as := getAS(id)
  if _, ok := e.camera[as]; !ok {
    e.camera[as] = []string{}
  }
  e.camera[as] = append(e.camera[as], id)
}

type byId []Candidate

func (a byId) Len() int           { return len(a) }
func (a byId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byId) Less(i, j int) bool { return a[i].Id < a[j].Id }

/* Run leader election in a specific as. */
func (e *Election) Elect(as string) []string {
  e.Lock()
  defer e.Unlock()

  candidates := e.candidates[as]

  // Sort the candidates by a criteria
  sort.Sort(byId(candidates))

  leaders := []string{}
  for i := 0; i < MaxLeaders; i++ {
    leaders = append(leaders, candidates[i].Id)
  }

  return leaders
}

func getNBR(id string) string {
  // We assume that ID is of form [AS].[NBR]
  return strings.Split(id, ".")[0]
}

func getAS(id string) string {
  // We assume that ID is of form [AS].[NBR]
  return strings.Split(id, ".")[0]
}
