package log

import (
  "runtime"
  "fmt"
)

var Log *Logger = NewLogger()

const (
  GetRedundancy = iota
  GetTime       = iota
  GetTraffic    = iota
  GetTimeCDF    = iota
  Stop          = iota
)

const maxTransfers int = 100000
const maxCompletes int = 1000

type piece struct {
  as  string
  idx int
}

type Logger struct {
  verbose    bool

  redundancy map[piece]int
  traffic    map[int]int
  times     []int

  transfers chan Transfer
  completes chan Completed
  queries   chan int

  stopped   bool
}

func NewLogger() *Logger {
  logger := &Logger{
    verbose : false,

    redundancy : make(map[piece]int),
    traffic    : make(map[int]int),
    times      : []int{},

    transfers  : make(chan Transfer, maxTransfers),
    completes  : make(chan Completed, maxCompletes),
    queries    : make(chan int, 1),

    stopped    : false,
  }

  go logger.run()

  return logger
}

/* Defaults*/
func SetVerbose(verbose bool)  { Log.SetVerbose(verbose) }
func Println(v ...interface{}) { Log.Println(v...) }
func LogCompleted(t Completed) { Log.LogCompleted(t) }
func LogTransfer(t Transfer)   { Log.LogTransfer(t) }
func Query(q int)              { Log.Query(q) }

/* Interface. */
func (l *Logger) SetVerbose(verbose bool) {
  l.verbose = verbose
}

func (l *Logger) Println(v ...interface{}) {
  if l.verbose {
    fmt.Println(v...)
  }
}

func (l *Logger) LogCompleted(t Completed) {
  l.completes <- t
}

func (l *Logger) LogTransfer(t Transfer) {
  l.transfers <- t
}

func (l *Logger) Query(q int) {
  l.queries <- q
}

/* Handlers. */
func (l *Logger) handleTransfer(t Transfer) {
  as := getAS(t.To)
  if as != getAS(t.From) {
    p := piece{
      as  : as,
      idx : t.Index,
    }
    if _, ok := l.redundancy[p]; !ok {
      l.redundancy[p] = 0
    }
    l.redundancy[p] += 1
  }

  idx := t.Index
  if _, ok := l.traffic[idx]; !ok {
    l.traffic[idx] = 0
  }
  l.traffic[idx] += 1
}

func (l *Logger) handleComplete(c Completed) {
  l.times = append(l.times, c.Time)
}

/* Queries. */
func (l *Logger) getRedundancy() {
  pieces := 0
  times  := 0
  for _, ctr := range l.redundancy {
    pieces += 1
    times  += ctr
  }
  redundancy := float64(times) / float64(pieces)
  fmt.Println("Redundancy:", redundancy)
}

func (l *Logger) getTraffic() {
  total := 0
  peers := 0
  for _, ctr := range l.traffic {
    total += ctr
    peers += 1
  }
  traffic := float64(total) / float64(peers)
  fmt.Println("Traffic:", traffic)
}

func (l *Logger) getTime() {
  fmt.Println("Average time:", getAverage(l.times))
  fmt.Println("50th percentile:", getPercentile(50.0, l.times))
  fmt.Println("90th percentile:", getPercentile(90.0, l.times))
}

func (l *Logger) getTimeCDF() {
  fmt.Print("Time CDF: [")
  for _, t := range normalize(l.times) {
    fmt.Print(t, ",")
  }
  fmt.Println("]")
}

/* Runner. */
func (l *Logger) run() {
  for {
    select {
    case t := <-l.transfers:
      l.handleTransfer(t)
      continue
    default:
    }

    select {
    case c := <-l.completes:
      l.handleComplete(c)
      continue
    default:
    }

    select {
    case q := <-l.queries:
      switch q {
      case GetRedundancy:
        l.getRedundancy()
      case GetTime:
        l.getTime()
      case GetTraffic:
        l.getTraffic()
      case GetTimeCDF:
        l.getTimeCDF()
      case Stop:
        l.stopped = true
      }
      continue
    default:
    }

    // All channels are drained
    if l.stopped {
      break
    }
    runtime.Gosched()
  }
}
