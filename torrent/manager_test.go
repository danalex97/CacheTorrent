package torrent

import (
  // "testing"
)

/* Mocks. */
type mockManager struct {
  uploads   []Upload
  downloads []*Download
}

func (m *mockManager) AddConnector(conn *Connector) {}
func (m *mockManager) Uploads()   []Upload    { return m.uploads }
func (m *mockManager) Downloads() []*Download { return m.downloads }

/* Tests. */
