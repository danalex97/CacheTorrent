package config

import (
  "reflect"
)

type Const interface {
  Value() int
}

type constant struct {
  field string
  value int
  init  bool
}

func NewConst(field string) Const {
  return &constant{
    field : field,
    init  : false,
  }
}

func (c *constant) Value() int {
  if !c.init {
    r := reflect.ValueOf(Config)
    f := reflect.Indirect(r).FieldByName(c.field)
    c.value = int(f.Int())
    c.init  = true
  }
  return c.value
}

const OutPeers string = "OutPeers"
const InPeers  string = "InPeers"

const MinNodes string = "MinNodes"
const Seeds    string = "Seeds"

const PieceSize string = "PieceSize"
const Pieces    string = "Pieces"

const Uploads     string = "Uploads"
const Optimistics string = "Optimistics"
const Interval    string = "Interval"

const Backlog string = "Backlog"
