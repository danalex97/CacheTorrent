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

  connector.Handshake  = NewHandshake(connector)
  connector.Upload     = NewUpload(connector)
  connector.Download   = NewDownload(connector)

  return connector
}

func (c *Connector) Run() {
  go c.Handshake.Run()
  go c.Upload.Run()
  go c.Download.Run()
}

func (c *Connector) Recv(m interface {}) {
  c.Handshake.Recv(m)
  c.Upload.Recv(m)
  c.Download.Recv(m)
}
