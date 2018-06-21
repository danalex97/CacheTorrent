## Remote runs

To be able to run simulations remotely you need:
- a network file system
- ability to log in using SSH

To run remote jobs, use the script `./remote.sh`:
```
Usage: ./remote.sh [-h] [-n NAME] [-r RUNS] [-k [KILL]]
                   [-pool POOL] [command [command ...]]

Run multiple simulations remotely.

Positional arguments:
  command

Optional arguments:
  -h, --help               Show this help message and exit.
  -n NAME, --name NAME     The name of the folder in which the results will be saved.
  -r RUNS, --runs RUNS     Number of times that the job runs.
  -k [KILL], --kill [KILL] Use this flag to kill all remote jobs.
  -pool POOL               Path to a .txt file containing IPs(or names) for the machine pool.

```

#### Example

Suppose that we want to simulate the following 3 commands:
```
go run main.go -ext=50
go run main.go -ext=40
go run main.go -ext=30
```

Let's say:
- the pool file is `pool.txt` in the project root folder
- we want to save the file in the folders
  - `results/ext30`
  - `results/ext40`
  - `results/ext50`
- we want to run each job 5 times and average the results

The remote run saves files in the folder `results` by default. We will, therefore, run:
```
./remote.sh --runs=5 --name=ext50 -pool=pool.txt go run main.go -ext=50
./remote.sh --runs=5 --name=ext40 -pool=pool.txt go run main.go -ext=40
./remote.sh --runs=5 --name=ext30 -pool=pool.txt go run main.go -ext=30
```

When we run the commands, we will be able to see the deployment log in file `results/ext50/log.txt`. Such a log looks like:
```
Job id: ext50
Running job: go run main.go -ext=50
Run remote job on: point20
Run remote job on: voxel01
Run remote job on: point37
Run remote job on: edge12
Run remote job on: arc02
Done: remote_run/point37.txt
Done: remote_run/point20.txt
Done: remote_run/edge12.txt
Done: remote_run/arc02.txt
Done: remote_run/voxel01.txt
Jobs finished. Starting callback.
Stopping server.
```

Jobs might fail or remain hanging for various reasons. (computer restarting, server crashing and so on) Suppose that the remote job `point20` is hanging. We can kill the job by using the command:
```
./remote.sh --kill point20
```

If we want to stop all the jobs in the network, we can run:
```
./remote.sh --kill
```

When all the simulation finish, a `summary.txt` file will be creates in the respective folder.

#### Pool files

A pool file is a file with a list of IPs or domains to which we will connect via SSH to deploy a simulation. An example of such a file is:
```
10.2.2.8
10.2.2.9
hero.com
alex.local
localhost
FE80::0202:B3FF:FE1E:8329
```

By default the script will use the Imperial Department of Computing pool of computers.

#### Experimental results

All our experimental results are saved into the archive `misc/results.zip`. Each folder contains multiple runs.

A single run is structured as follows:
- log.txt
- runs
  - 0.txt
  - 1.txt
  - ...
  - computer1.txt
  - computer2.txt
  - ...
- summary.txt

In all *numbered files* there is only the individual summary of each run, containing the metrics:
- 50th download time percentile
- 90th download time percentile
- average download time
- number of redundant transmissions per piece

The *computer.txt* files contain the full logs of each separate run. The *summary.txt* file contains the averages of all the metrics presented above. Some *summary.txt* files will also contains the leader and follower 50th percentile, 90th percentile and average download time. In some *computer.txt* at the end CDFs can be found if the query that run contained the `-cdf` flag.
