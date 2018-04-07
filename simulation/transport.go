package simulation

/* All the interfaces as provided by the Speer simulator. */
type Data struct {
  Id   string
  Size int
}

type Node interface {
  Up()   int
  Down() int
}

type Link interface {
  Upload(Data)
  Download() <-chan Data

  From() Node
  To()   Node
}

/* Internal interface identical to Speer's Engine interface. */
type Transport interface {
  Node

  Connect(string) Link

  ControlPing(string) bool
  ControlSend(string, interface {})
  ControlRecv() <-chan interface {}
}
