package log

import (
  "runtime"
)

var Log *Logger = NewLogger()

const (
  GetRedundancy = iota
  Stop          = iota
)

const maxTransfers int = 100000

type piece struct {
  as  string
  idx int
}

type Logger struct {
  redundancy map[piece]int

  transfers chan Transfer
  queries   chan int

  stopped   bool
}

func NewLogger() *Logger {
  logger := &Logger{
    redundancy : make(map[piece]int),

    transfers  : make(chan Transfer, maxTransfers),
    queries    : make(chan int),

    stopped    : false,
  }

  go logger.run()

  return logger
}

func (l *Logger) LogTransfer(t Transfer) {
  l.transfers <- t
}

func (l *Logger) Query(q int) {
  l.queries <- q
}

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
}

func (l *Logger) getRedundancy() {

}

func (l *Logger) run() {
  for {
    select {
    case t := <-l.transfers:
      l.handleTransfer(t)
    default:
      continue
    }

    select {
    case q := <-l.queries:
      switch q {
      case GetRedundancy:
        l.getRedundancy()
      case Stop:
        l.stopped = true
      }
    default:
      continue
    }

    // All channels are drained
    if l.stopped {
      break
    }
    runtime.Gosched()
  }
}
