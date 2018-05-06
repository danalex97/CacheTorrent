# nfsTorrent
[![Build Status](https://travis-ci.org/danalex97/Speer.svg?branch=master)](https://travis-ci.org/danalex97/nfsTorrent) [![Coverage Status](https://coveralls.io/repos/github/danalex97/nfsTorrent/badge.svg?branch=master)](https://coveralls.io/github/danalex97/nfsTorrent?branch=master)

Network-friendly speedy torrents.

To run use `go run main.go`. It supports the following arguments:
```
-bias int
  	Number of outgoing connections for a biased Tracker.
-cdf
  	Enable printing time cumulative distribution function.
-conf string
  	The path to configuration .json file. (default "./confs/small.json")
-cpuprofile file
  	Write cpu profile to file.
-ext int
  	Use the textesion with ext percent number of leaders.
-memprofile file
  	Write memory profile to file.
-v	Verbose output
```

To run multiple simulations in parallel, use `./remote.sh`.
```
Usage: ./remote.sh [-h] [-n NAME] [-r RUNS] command

Run multiple simulations remotely.

Positional arguments:
  command

Optional arguments:
  -h, --help            show this help message and exit
  -n NAME, --name NAME  The name of the folder in which the results will be
                        saved.
  -r RUNS, --runs RUNS  Number of times that the job runs.
```
Example usage:
```
./remote.sh -n=test -r=5 go run main.go -conf=confs/itl.json -v
```


### Simulation [![GoDoc](https://godoc.org/github.com/danalex97/Speer/interfaces?status.png)](https://godoc.org/github.com/danalex97/Speer/interfaces)

The simulation package uses [Speer](https://github.com/danalex97/Speer) simulator.

### BitTorrent [![GoDoc](https://godoc.org/github.com/danalex97/nfsTorrent/torrent?status.png)](https://godoc.org/github.com/danalex97/nfsTorrent/torrent)

A simplified implementation of the BitTorrent protocol.

### CacheTorrent [![GoDoc](https://godoc.org/github.com/danalex97/nfsTorrent/cache_torrent?status.png)](https://godoc.org/github.com/danalex97/nfsTorrent/cache_torrent)

A BitTorrent extension aimed to make the protocol more network friendly using leader election, caches and indirect requests.
