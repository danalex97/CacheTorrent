package torrent

// The Connector is only an interface towards Upload and Download components.
// For simplicity at building the Connector is initialized using a builder
// pattern(with no build intermediary structure).
type Connector struct {
  *Components

  From string
  To   string

  Upload    Upload
  Download  Download
  Handshake Handshake

  runHandshake bool
}

func NewConnector(from, to string, components *Components) *Connector {
  connector := new(Connector)

  connector.Components = components

  connector.From      = from
  connector.To        = to

  connector.Handshake = NewHandshake(connector)

  return connector
}

// Decorate the connection with a new Upload. To support multiple types of
// download, we use a 'func(*Connector) Download' for interfacing.
func (c *Connector) WithUpload(newUpload func(*Connector) Upload) *Connector {
  c.Upload = newUpload(c)
  return c
}

// Decorate the connection with a new Download. To support multiple types of
// download, we use a 'func(*Connector) Download' for interfacing.
func (c *Connector) WithDownload(newDownload func(*Connector) Download) *Connector {
  c.Download = newDownload(c)
  return c
}

// The Register method is used to bind the Connection to a Peer. The Peer will
// call this method from `AddConnector(id string)`.
func (c *Connector) Register(p *Peer) *Connector {
  p.Connectors[c.To] = c
  p.Manager.AddConnector(c)

  go c.Run()
  return c
}

// The Connector is a Runner, thus having the Run and Recv(m interface {})
// methods. The Run method should be called when the connector is registered.
func (c *Connector) Run() {
  if c.Upload != nil {
    go c.Upload.Run()
    go c.Handshake.Run()
  }
  if c.Download != nil {
    go c.Download.Run()
  }
}

// The Recv method is called whenever a message is dispached for processing
// towards the Connector.
func (c *Connector) Recv(m interface {}) {
  if c.Upload != nil {
    c.Upload.Recv(m)
  }
  if c.Download != nil {
    c.Download.Recv(m)
    c.Handshake.Recv(m)
  }
}
