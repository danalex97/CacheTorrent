package log

type Transfer struct {
  From string
  To   string

  Index  int
}

type Completed struct {
  Time int
}
