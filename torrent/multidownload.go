package torrent

type MultiDownload struct {
}

func NewMultiDownload() *MultiDownload {
  return &MultiDownload{}
}

func (d *MultiDownload) Lost(lost []int) {
  // handle pieces that were lost due to chocking
}
