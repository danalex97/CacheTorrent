package torrent

import (
)

/* Mocks. */
type mockRunner struct {}

func (m *mockRunner) Run() {}
func (m *mockRunner) Recv(interface {}) {}
