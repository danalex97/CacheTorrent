package torrent

/**
 * This file follows the 'Upload' BitTorrent 5.3.0 release.
 *
 * The connector is only an interface towards upload and download.
 */
type Connector struct {
  *Components

  From string
  To   string

  Upload    Upload
  Download  Download
  Handshake Handshake
}

func NewConnector(from, to string, components *Components) *Connector {
  connector := new(Connector)

  connector.Components = components

  connector.From = from
  connector.To   = to

  return connector
}

func (c *Connector) WithHandshake() *Connector {
  c.Handshake = NewHandshake(c)
  return c
}

func (c *Connector) WithUpload() *Connector {
  c.Upload = NewUpload(c)
  return c
}

func (c *Connector) WithDownload() *Connector {
  c.Download = NewDownload(c)
  return c
}

func (c *Connector) Register(p *Peer) *Connector {
  p.Connectors[c.To] = c
  p.Manager.AddConnector(c)

  go c.Run()
  return c
}

func (c *Connector) Run() {
  if c.Handshake != nil {
    go c.Handshake.Run()
  }
  if c.Upload != nil {
    go c.Upload.Run()
  }
  if c.Download != nil {
    go c.Download.Run()
  }
}

func (c *Connector) Recv(m interface {}) {
  if c.Handshake != nil {
    c.Handshake.Recv(m)
  }
  if c.Upload != nil {
    c.Upload.Recv(m)
  }
  if c.Download != nil {
    c.Download.Recv(m)
  }
}
