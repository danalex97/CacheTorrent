package cache_torrent

import (
  . "github.com/danalex97/Speer/interfaces"
  "sync"
  "sort"
)

const MaxLeaders int = 1

type Election struct {
  sync.Mutex

  limit       int
  nodes       int

  camera      map[string][]string
  candidates  map[string][]Candidate

  elected     map[string][]string
  transport   Transport
}

func NewElection(limit int, transport Transport) *Election {
  return &Election{
    camera     : make(map[string][]string),
    candidates : make(map[string][]Candidate),
    elected    : make(map[string][]string),
    limit      : limit,
    nodes      : 0,
    transport  : transport,
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

  e.nodes++

  /* Add candidate to candidate list. */
  as := getAS(candidate.Id)
  if _, ok := e.candidates[as]; !ok {
    e.candidates[as] = []Candidate{}
  }
  e.candidates[as] = append(e.candidates[as], candidate)

  // When we reach the node limit, we run the full elections.
  if e.nodes == e.limit {
    e.Unlock()
    e.RunElection()
    e.Lock()
  }

  // For ulterior joins, we only respond with the Leader messages.
  if e.nodes > e.limit {
    elected, ok := e.elected[as]
    if !ok {
      // If there are no leaders, we designate the requester as a leader.
      // i.e. the node will, thus, follow the original BitTorrent protocol
      elected = []string{candidate.Id}
    }
    e.transport.ControlSend(candidate.Id, Leaders{elected})
  }
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

/* Run elections for all ASes. */
func (e *Election) RunElection() {
  e.Lock()
  defer e.Unlock()

  for as, _ := range e.camera {
    e.Unlock()
    e.elected[as] = e.Elect(as)
    e.Lock()
  }

  // Send Leader messages
  for as, camera := range e.camera {
    elected := e.elected[as]
    for _, node := range camera {
      e.transport.ControlSend(node, Leaders{elected})
    }
  }
}

type byId []Candidate

func (a byId) Len() int           { return len(a) }
func (a byId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byId) Less(i, j int) bool { return a[i].Id < a[j].Id }

/* Run leader election in a specific as. */
func (e *Election) Elect(as string) []string {
  e.Lock()
  defer e.Unlock()

  candidates, ok := e.candidates[as]
  if !ok {
    // If there are no candidates, we designate all nodes as leaders,
    // i.e. each node will be able to communicate with the exterior
    return e.camera[as]
  }

  // Sort the candidates by a criteria
  sort.Sort(byId(candidates))

  leaders    := []string{}
  maxLeaders := MaxLeaders
  if len(candidates) < maxLeaders {
    maxLeaders = len(candidates)
  }
  for i := 0; i < maxLeaders; i++ {
    leaders = append(leaders, candidates[i].Id)
  }

  return leaders
}
