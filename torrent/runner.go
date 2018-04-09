package torrent

type Runner interface {
  Run()
  Recv(interface {})
}
