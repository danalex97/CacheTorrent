package torrent

import (
  // "testing"
)

/* Mocks. */
type mockUpload struct {
  mockRunner

  isInterested bool
  choke        bool
  rate         float64
}

func NewMockUpload(rate float64) Upload {
  return &mockUpload{
    isInterested : false,
    choke : true,
    rate : rate,
  }
}

func (m *mockUpload) Choke() {
  m.choke = true
}

func (m *mockUpload) Unchoke() {
  m.choke = false
}

func (m *mockUpload) Choking() bool {
  return m.choke
}

func (m *mockUpload) IsInterested() bool {
  return m.isInterested
}

func (m *mockUpload) Rate() float64 {
  return m.rate
}

/* Tests. */
