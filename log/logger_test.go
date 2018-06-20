package log

import (
  "testing"
  "fmt"
)

func initTest(out chan<- string) {
  printf = func(format string, args ...interface {}) (int, error) {
    out <- fmt.Sprintf(format, args)
    return 0, nil
  }
  println = func(args ...interface {}) (int, error) {
    out <- fmt.Sprint(args)
    return 0, nil
  }
  print = println
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
    t.Fatalf("%s != %s", a, b)
	}
}

func TestLoggerGetsCorrectMetricsBitTorrent(t *testing.T) {
  l := NewLogger()

  out := make(chan string, 1000)
  initTest(out)

  l.LogCompleted(Completed{"1.1", 1000})
  l.LogCompleted(Completed{"1.2", 2000})
  l.LogCompleted(Completed{"2.1", 3000})
  l.LogCompleted(Completed{"2.2", 4000})
  l.LogCompleted(Completed{"2.3", 4000})

  l.LogTransfer(Transfer{"1.1", "1.2", 1})
  l.LogTransfer(Transfer{"1.1", "2.1", 1})
  l.LogTransfer(Transfer{"1.1", "2.2", 1})
  l.LogTransfer(Transfer{"1.1", "2.3", 1})

  l.Query(GetRedundancy)
  l.Query(GetTime)
  l.Query(GetTraffic)
  l.Query(GetTimeCDF)

  l.Query(Stop)

  assertEqual(t, <-out, "[Redundancy: 3]")
  assertEqual(t, <-out, "[Average time: 1800]")
  assertEqual(t, <-out, "[50th percentile: 2000]")
  assertEqual(t, <-out, "[90th percentile: 3000]")
  assertEqual(t, <-out, "[Traffic: 4]")

  <-out
  assertEqual(t, (<-out), "[0 ,]")
  assertEqual(t, (<-out), "[1000 ,]")
  assertEqual(t, (<-out), "[2000 ,]")
  assertEqual(t, (<-out), "[3000 ,]")
  assertEqual(t, (<-out), "[3000 ,]")
}

func TestLoggerGetsCorrectMetricsCacheTorrent(t *testing.T) {
  l := NewLogger()

  out := make(chan string, 1000)
  initTest(out)

  l.LogLeader(Leader{"1.1"})
  l.LogLeader(Leader{"2.1"})

  l.LogCompleted(Completed{"1.1", 1000})
  l.LogCompleted(Completed{"1.2", 2000})
  l.LogCompleted(Completed{"2.1", 3000})
  l.LogCompleted(Completed{"2.2", 4000})
  l.LogCompleted(Completed{"2.3", 4000})

  l.LogTransfer(Transfer{"1.1", "1.2", 1})
  l.LogTransfer(Transfer{"1.1", "2.1", 1})
  l.LogTransfer(Transfer{"1.1", "2.2", 1})
  l.LogTransfer(Transfer{"1.1", "2.3", 1})

  l.Query(GetRedundancy)
  l.Query(GetTime)
  l.Query(GetTraffic)
  l.Query(GetTimeCDF)
  l.Query(GetTimeLeader)
  l.Query(GetLeaderCDF)

  l.Query(Stop)

  assertEqual(t, <-out, "[Redundancy: 3]")
  assertEqual(t, <-out, "[Average time: 1800]")
  assertEqual(t, <-out, "[50th percentile: 2000]")
  assertEqual(t, <-out, "[90th percentile: 3000]")
  assertEqual(t, <-out, "[Traffic: 4]")

  <-out
  assertEqual(t, (<-out), "[0 ,]")
  assertEqual(t, (<-out), "[1000 ,]")
  assertEqual(t, (<-out), "[2000 ,]")
  assertEqual(t, (<-out), "[3000 ,]")
  assertEqual(t, (<-out), "[3000 ,]")
  <-out
  assertEqual(t, (<-out), "[Leader 50th percentile: 2000]")
  assertEqual(t, (<-out), "[Leader 90th percentile: 2000]")
  assertEqual(t, (<-out), "[Follower 50th percentile: 3000]")
  assertEqual(t, (<-out), "[Follower 90th percentile: 3000]")

  <-out
  assertEqual(t, (<-out), "[0 ,]")
  assertEqual(t, (<-out), "[2000 ,]")
  <-out

  <-out
  assertEqual(t, (<-out), "[0 ,]")
  assertEqual(t, (<-out), "[1000 ,]")
  assertEqual(t, (<-out), "[3000 ,]")
  <-out
}
