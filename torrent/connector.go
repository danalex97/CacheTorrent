package torrent

// This file follows the 'Upload' BitTorrent 5.3.0 release

type Connector struct {
  from string
  to   string

  upload    *Upload
  download  *Download
  handshake *Handshake

  components *Components
}

func NewConnector(from, to string, components *Components) *Connector {
  connector := new(Connector)

  connector.from  = from
  connector.to    = to

  connector.components = components
  connector.handshake  = NewHandshake(connector)
  connector.upload     = NewUpload(connector)
  connector.download   = NewDownload(connector)

  return connector
}

func (c *Connector) Run() {
  go c.handshake.Run()
  go c.upload.Run()
  go c.download.Run()
}

func (c *Connector) Recv(m interface {}) {
  c.handshake.Recv(m)
  c.upload.Recv(m)
  c.download.Recv(m)
}
