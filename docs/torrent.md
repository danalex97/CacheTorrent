## BitTorrent [![GoDoc](https://godoc.org/github.com/danalex97/nfsTorrent/torrent?status.png)](https://godoc.org/github.com/danalex97/nfsTorrent/torrent)

The BitTorrent package is a simplified implementation of the BitTorrent protocol. For individual component implementation please refer to the **godoc**.

### General protocol

**BitTorrent** is a communication protocol for peer-to-peer file sharing (P2P) which is used to distribute data and electronic files over the Internet. The BitTorrent overlay is an unstructured one. Each peer wants to download the same file and the peers collaborate with the purpose of downloading the file faster and distributing the bandwidth load.

A **tracker** is a special node which has the responsibility of providing the set of neighbors for each peer. The other peers are called **leechers**. The file to be transferred is broken into pieces, which are shared between the peers. A **leecher** that contains all the pieces and uploads them to other peers is called a seed.

### nfsTorrent/torrent

We implement a version of the BitTorrent protocol following the [version 5.3 design](https://github.com/danalex97/BitTorrent) with some modifications. The main components in our implementation are:
- **Storage** - stores, checks and responds to queries about the piece downloads
- **Picker** - for a Download connection, chooses the next piece to make a request for
- **Choker** - chooses peers to choke and unchoke, sending sends the corresponding
messages
- **Connection Manager** - updates the list of alive connections

Each connection is represented by a **Connector** structure. The **Connector** is composed of:
- **Upload** - starts uploads when it receives a request for a piece and notifies rest of components when the upload is choked
- **Download** - reacts to most types of messages, establishing which pieces should be downloaded by using the Picker
- **Handshake** - establishes a bidirectional or unidirectional connection used for data transfer

### Message Types

We follow the original BitTorrent protocol for most message types:
- Choke
- Unchoke
- Interested
- NotInterested
- Have
- Request
- Piece

We modify some messages and we provide new messages for integration with the [Speer](https://github.com/danalex97/Speer) simulator:
- Join
- Neighbours
- TrackerReq
- TrackerRes
- SeedReq
- SeedRes
- ConnReq

We eliminate the message following message types:
- Cancel

For individual message descriptions please consult the [documentation](https://godoc.org/github.com/danalex97/nfsTorrent/torrent).
