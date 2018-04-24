package log

import (
  "strings"
  "sort"
)

func getAS(id string) string {
  // We assume that ID is of form [AS].[NBR]
  return strings.Split(id, ".")[0]
}

func getPercentile(percentile float64, data []int) float64 {
  norm := normalize(data)
  idx := int(float64(len(norm)) * percentile / 100.0)

  if idx >= len(norm) {
    idx = len(norm) - 1
  }

  return float64(norm[idx])
}

func getAverage(data []int) float64 {
  norm := normalize(data)

  var sum int64
  var ctr int
  for _, v := range norm {
    sum += int64(v)
    ctr += 1
  }
  return float64(sum) / float64(ctr)
}

func normalize(data []int) []int {
  mn := data[0]
  for _, v := range data {
    if v < mn {
      mn = v
    }
  }

  norm := []int{}
  for _, v := range data {
    norm = append(norm, v - mn)
  }

  sort.Ints(norm)

  return norm
}
