package config

import (
  "sync"
)

type WGProgress struct {
  wg sync.WaitGroup

  size int
}

func NewWGProgress() *WGProgress {
  return &WGProgress{
    size : 0,
  }
}

func (p *WGProgress) Add() {
  p.size++
}

func (p *WGProgress) Progress() {
  p.wg.Done()
}

func (p *WGProgress) Advance() {
  p.wg.Add(p.size)
  p.wg.Wait()
}
