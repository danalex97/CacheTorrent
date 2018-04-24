package log

import (
  "strings"
)

func getAS(id string) string {
  // We assume that ID is of form [AS].[NBR]
  return strings.Split(id, ".")[0]
}

func getPercentile(percentile float64, data []int) float64 {
  return 0
}

func getAverage(data []int) float64 {
  var sum int64
  var ctr int
  for _, v := range data {
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

  return norm
}
