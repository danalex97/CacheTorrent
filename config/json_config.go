package config

import (
  "encoding/json"
  "io/ioutil"
  "strings"
  "bytes"
)

func trimComments(arr []byte) []byte {
  str := string(bytes.Trim(arr, "\x00"))

  out := ""
  for _, v := range strings.Split(str, "\n") {
    if len(strings.TrimSpace(v)) > 3 {
        if strings.TrimSpace(v)[:2] != "//" {
            out = out + v
        }
    } else {
        out = out + v
    }
    out += "\n"
  }

  return []byte(out)
}

func JSONConfig(path string) *Conf {
  raw, err := ioutil.ReadFile(path)
  if err != nil {
    panic(err.Error())
  }

  raw = trimComments(raw)

  var conf Conf
  json.Unmarshal(raw, &conf)

  conf.SharedCallback = func () {}
  conf.SharedInit     = func () {}

  return &conf
}
