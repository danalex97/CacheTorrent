## CacheTorrent [![GoDoc](https://godoc.org/github.com/danalex97/nfsTorrent/cache_torrent?status.png)](https://godoc.org/github.com/danalex97/nfsTorrent/cache_torrent)

![ ](pics/cache.png)

The general idea of the **CacheTorrent** protocol is assigning the role of leaders to a percent of the peers. Each Leader will act as a unidirectional cache for downloading in each autonmous system. On one hand, a Follower will have to connect to a Leader in order to download a piece. On the other hand, a Follower can upload directly to another Leader in a different autonmous system.

### Implementation

We extend the BitTorrent protocol by adding components on top of it. We decorate a BitTorrent **Peer** as follows:
- **Peer** - a CacheTorrent wrapper over a BitTorrent peer

### Message Types

Besides the [BitTorrent message types](torrent.md), we have the following new messages:
- LeaderStart - message used to establish an indirect connection(connection to another peer through a Leader)
- Neighbours - similar to Neighbours message in BitTorrent protocol
- Candidate - message used by a peer to become a candidate in leader election
- Leader - message sent by the Tracker to announce the leaders
