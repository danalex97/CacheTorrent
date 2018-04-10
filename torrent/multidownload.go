package torrent

type MultiDownload struct {
}

func NewMultiDownload() *MultiDownload {
  return &MultiDownload{}
}

func (d *MultiDownload) Lost(lost []pieceMeta) {
  // handle pieces that were lost due to chocking
}
