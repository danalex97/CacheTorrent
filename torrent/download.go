package torrent

import (
  // "fmt"
)

// This file follows the 'download' BitTorrent 5.3.0 release

type Download struct {
  *Components

  me   string // the node that downloads
  from string // the node that we download from

  connector *Connector
}

func NewDownload(connector *Connector) Runner {
  return &Download{
    connector.components,

    connector.from,
    connector.to,

    connector,
  }
}

func (d *Download) Run() {

}

func (d *Download) Recv(m interface {}) {
  switch msg := m.(type) {
  case have:
    index := msg.index

    // send interested if I'm not interested and chocked
    if d.connector.chocked && !d.connector.interested {
      // I need to be interested in the piece as well
      if _, ok := d.Storage.Have(index); !ok {
        // Send interested message to node
        d.Transport.ControlSend(d.from, interested{d.me})
      }
    }

    // let picker know I can get piece index
    d.Picker.GotHave(index)
  }
}
