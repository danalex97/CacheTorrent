package cache_torrent

import (
  "github.com/danalex97/nfsTorrent/torrent"
)

func NewLeaderLocalConnector(from, to string, components *torrent.Components) *torrent.Connector {
  connector := new(torrent.Connector)

  connector.Components = components

  connector.From  = from
  connector.To    = to

  connector.Handshake  = torrent.NewHandshake(connector)
  connector.Upload     = torrent.NewUpload(connector)
  connector.Download   = torrent.NewDownload(connector)

  return connector
}

func NewLeaderPeerConnector(from, to string, components *torrent.Components) *torrent.Connector {
  connector := new(torrent.Connector)

  connector.Components = components

  connector.From  = from
  connector.To    = to

  connector.Handshake  = torrent.NewHandshake(connector)
  connector.Upload     = torrent.NewUpload(connector)
  connector.Download   = torrent.NewDownload(connector)

  return connector
}
