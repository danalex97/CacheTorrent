package torrent

type pieceMeta struct {
  index  int
  begin  int
  length int
}

type Components struct {
  Picker     *Picker
  Storage    *Storage
}
