package torrent

import (
  "testing"
)

/* Mocks. */
type mockConst struct {
  value interface {}
}

func (c *mockConst) Ref() interface {} { return c.value }
func (c *mockConst) Int() int          { return c.value.(int) }

type mockChoker struct {
  interestedCalled    bool
  notInterestedCalled bool
}

func (c *mockChoker) Interested(conn Upload) { c.interestedCalled = true }
func (c *mockChoker) NotInterested(conn Upload) { c.notInterestedCalled = true }
func (c *mockChoker) Run() {}

/* Testing. */
func makeChoker(uploads []Upload) ([]Upload, *choker) {
  return uploads, NewChoker(&mockManager{
    uploads : uploads,
  }, func() int {return 0}, false).(*choker)
}

func TestRechokingDecreasingByRate(t *testing.T) {
  uploads     = &mockConst{3}
  optimistics = &mockConst{0}

  ups, c := makeChoker([]Upload{
    &mockUpload{isInterested:true, choke: true, rate:10},
    &mockUpload{isInterested:true, choke: true, rate:40},
    &mockUpload{isInterested:true, choke: true, rate:20},
    &mockUpload{isInterested:true, choke: true, rate:30},
    &mockUpload{isInterested:true, choke: true, rate:5},
  })

  c.rechoke()

  assertEqual(t, ups[0].Choking(), true)
  assertEqual(t, ups[1].Choking(), false)
  assertEqual(t, ups[2].Choking(), false)
  assertEqual(t, ups[3].Choking(), false)
  assertEqual(t, ups[4].Choking(), true)
}

func TestRechokingOptimisticChosenOutOfRemaining(t *testing.T) {
  uploads     = &mockConst{3}
  optimistics = &mockConst{1}

  ups, c := makeChoker([]Upload{
    &mockUpload{isInterested:true, choke: true, rate:10},
    &mockUpload{isInterested:true, choke: true, rate:40},
    &mockUpload{isInterested:true, choke: true, rate:20},
    &mockUpload{isInterested:true, choke: true, rate:30},
    &mockUpload{isInterested:true, choke: true, rate:5},
  })

  c.rechoke()

  assertEqual(t, ups[1].Choking(), false)
  assertEqual(t, ups[2].Choking(), false)
  assertEqual(t, ups[3].Choking(), false)
  assertEqual(t, ups[4].Choking() || ups[0].Choking(), true)
  assertEqual(t, ups[0].Choking() && ups[0].Choking(), false)
}

func TestRechokingDiffrentOptimistics(t *testing.T) {
  uploads     = &mockConst{0}
  optimistics = &mockConst{1}

  idxs := []int{}

  for i := 0; i < 10; i++ {
    ups, c := makeChoker([]Upload{
      &mockUpload{isInterested:true, choke: true, rate:10},
      &mockUpload{isInterested:true, choke: true, rate:40},
      &mockUpload{isInterested:true, choke: true, rate:20},
      &mockUpload{isInterested:true, choke: true, rate:30},
      &mockUpload{isInterested:true, choke: true, rate:5},
    })

    c.rechoke()

    for j := 0; j < 5; j++ {
      if !ups[j].Choking() {
        idxs = append(idxs, j)
      }
    }
  }

  for i := 0; i < 10; i++ {
    for j := 0; j < 10; j++ {
      if idxs[i] != idxs[j] {
        return
      }
    }
  }

  t.Fatalf("Keep choosing same optimistics.")
}

func TestRechokingIgnoresNotIterestedConnections(t *testing.T) {
  uploads     = &mockConst{3}
  optimistics = &mockConst{1}

  ups, c := makeChoker([]Upload{
    &mockUpload{isInterested:true,  choke: true, rate:10},
    &mockUpload{isInterested:false, choke: true, rate:40},
    &mockUpload{isInterested:false, choke: true, rate:20},
    &mockUpload{isInterested:true,  choke: true, rate:30},
    &mockUpload{isInterested:true,  choke: true, rate:5},
  })

  c.rechoke()

  assertEqual(t, ups[0].Choking(), false)
  assertEqual(t, ups[1].Choking(), true)
  assertEqual(t, ups[2].Choking(), true)
  assertEqual(t, ups[3].Choking(), false)
  assertEqual(t, ups[4].Choking(), false)
}

func TestRechokeCalled(t *testing.T) {
  uploads = &mockConst{0}
  optimistics = &mockConst{0}

  ups, c := makeChoker([]Upload{
    &mockUpload{isInterested:true, choke: false, rate:10},
  })
  c.Interested(&mockUpload{choke: false})
  assertEqual(t, ups[0].Choking(), true)

  ups, c = makeChoker([]Upload{
    &mockUpload{isInterested:true, choke: false, rate:10},
  })
  c.NotInterested(&mockUpload{choke: false})
  assertEqual(t, ups[0].Choking(), true)
}
