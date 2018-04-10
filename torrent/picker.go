package torrent

type Picker struct {
  pieces map[int]pieceMeta // the pieces that I have
}

func NewPicker(pieces []pieceMeta) *Picker {
  picker := new(Picker)
  picker.pieces = make(map[int]pieceMeta)

  for _, p := range pieces {
    picker.pieces[p.index] = p
  }
  return picker
}

func (p *Picker) GotHave(index int) {
}
