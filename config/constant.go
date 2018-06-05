package config

import (
  "reflect"
  "sync"
)

type Const interface {
  Ref() interface {}
  Int()  int
}

type constant struct {
  sync.RWMutex

  field string
  value interface {}
  init  bool
}

/**
 * Usage:
 *  var c config.Const = config.NewConst(config.OutPeers)
 */
func NewConst(field string) Const {
  return &constant{
    field : field,
    init  : false,
  }
}

func (c *constant) Ref() interface {} {
  c.RLock()
  defer c.RUnlock()

  if !c.init {
    c.RUnlock()
    c.Lock()

    r := reflect.ValueOf(Config)
    f := reflect.Indirect(r).FieldByName(c.field)
    c.value = f.Interface()
    c.init  = true

    c.Unlock()
    c.RLock()
  }
  return c.value
}

func (c *constant) Int() int {
  return c.Ref().(int)
}

/**
 * Constant values to pass to "config.NewConst".
 */
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

const LeaderPercent string = "LeaderPercent"
const Bias string = "Bias"

const AllNodesRun string = "AllNodesRun"

const Multi string = "Multi"
const StoragePieces string = "StoragePieces"
