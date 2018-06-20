# Network friendly speedy Torrents
[![Build Status](https://travis-ci.org/danalex97/Speer.svg?branch=master)](https://travis-ci.org/danalex97/nfsTorrent) [![Coverage Status](https://coveralls.io/repos/github/danalex97/nfsTorrent/badge.svg?branch=master)](https://coveralls.io/github/danalex97/nfsTorrent?branch=master)

CacheTorrent is a [file sharing system](https://en.wikipedia.org/wiki/File_sharing) CacheTorrent based on leader election, caches and indirect requests. It is fast in homogeneous networks and inexpensive in terms of inter-ISP redundant transmissions under any network conditions.

The CacheTorrent system extends the [BitTorrent protocol](https://en.wikipedia.org/wiki/BitTorrent), providing strong incentives for users to follow the proposed solution.

### Quick Start

- Prerequisites: [Golang](https://golang.org/) >= 1.6
- [Installation](docs/install.md) guide
- [Usage](docs/usage.md) instructions
- [Testing](docs/testing.md)
- Reading [experimental results](docs/results.md)

### Documentation

The general documentation is provided in Markdown, while the code uses [GoDoc](https://godoc.org/).

#### Protocols
  - **[BitTorrent](docs/torrent.md)** - a simplified implementation of the BitTorrent protocol
  - **[CacheTorrent](docs/cache.md)** - a BitTorrent extension aimed to make the protocol mode network friendly by reducing the inter-ISP traffic
  - **[Extensions](docs/extension.md)** - extensions on top of CacheTorrent based on emperical results

#### Tools
  - Visualization tool [[usage](docs/usage.md)]
  - Remote deployment tool [[usage](docs/usage.md), [experimental results](docs/results.md)]

Simulations are made using the [Speer](https://github.com/danalex97/Speer) simulator.

### How to contribute!
Want to help? You can [raise an issues](https://help.github.com/articles/creating-an-issue/) or contact me directly at *dan.alex97@yahoo.com*.
