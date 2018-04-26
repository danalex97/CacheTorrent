# nfsTorrent
[![Build Status](https://travis-ci.org/danalex97/Speer.svg?branch=master)](https://travis-ci.org/danalex97/nfsTorrent) [![Coverage Status](https://coveralls.io/repos/github/danalex97/nfsTorrent/badge.svg?branch=master)](https://coveralls.io/github/danalex97/nfsTorrent?branch=master)

Network-friendly speedy torrents.

To run use `go run main.go`. It supports the following arguments:
```
-conf string
    The path to configuration .json file. (default "./confs/small.json")
-ext
    Whether we use the extension
-v
    Verbose output
```

To run multiple simulations in parallel, use `python3 remote.py`. The jobs can
be configured by hand as follows:
```python
jobs = [
    Job(pool, "go run main.go -conf=confs/small.json", times=10),
    Job(pool, "go run main.go -ext -conf=confs/small.json", times=10),
]
```


### Simulation [![GoDoc](https://godoc.org/github.com/danalex97/Speer/interfaces?status.png)](https://godoc.org/github.com/danalex97/Speer/interfaces)

The simulation package uses [Speer](https://github.com/danalex97/Speer) simulator.

### BitTorrent [![GoDoc](https://godoc.org/github.com/danalex97/nfsTorrent/torrent?status.png)](https://godoc.org/github.com/danalex97/nfsTorrent/torrent)

A simplified implementation of the BitTorrent protocol.

### CacheTorrent [![GoDoc](https://godoc.org/github.com/danalex97/nfsTorrent/cache_torrent?status.png)](https://godoc.org/github.com/danalex97/nfsTorrent/cache_torrent)

A BitTorrent extension aimed to make the protocol more network friendly using leader election, caches and indirect requests.
