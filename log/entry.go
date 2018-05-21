package log

type Transfer struct {
  From string
  To   string

  Index  int
}

type Completed struct {
  Id   string
  Time int
}

type Leader struct {
  Id string
}

type Packet struct {
  Src  string `json:"src"`
  Dst  string `json:"dst"`
  Type string `json:"type"`
}
