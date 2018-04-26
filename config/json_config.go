package config

import (
  "encoding/json"
  "io/ioutil"
  "fmt"
)

func JSONConfig(path string) *Conf {
  raw, err := ioutil.ReadFile(path)
  if err != nil {
    panic(err.Error())
  }

  var conf Conf
  json.Unmarshal(raw, &conf)

  conf.SharedCallback = func () {}
  conf.SharedInit     = func () {}

  fmt.Println(conf)

  return &conf
}
