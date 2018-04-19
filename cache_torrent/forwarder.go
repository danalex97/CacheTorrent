package cache_torrent

type Forwarder struct {
  from string
  to   string
}

func NewForwarder(from, to string) *Forwarder {
  return &Forwarder{
    from : from,
    to   : to,
  }
}

func (f *Forwarder) Recv(m interface {}) {
}
