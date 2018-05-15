import os
import subprocess
import socket
import random
import string

def get_ip():
    """Bind the process to a socket."""
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect(("8.8.8.8", 80))

    ip = s.getsockname()[0]
    s.close()

    return ip

ID  = "ad5915"
EXT = "doc.ic.ac.uk"

KEYS = {
  "red"  : "Redundancy",
  "time" : "Average time",
  "50p"  : "50th percentile",
  "90p"  : "90th percentile",
  "l50p" : "Leader 50th percentile",
  "l90p" : "Leader 90th percentile",
  "f50p" : "Follower 50th percentile",
  "f90p" : "Follower 90th percentile",
}

random_id = lambda : ''.join(
    random.choice(string.ascii_lowercase + string.digits)
    for _ in range(20))

def test_remote(id, host):
    # Check if there is only one user and
    # my script is not already running on the host.
    print("Test remote {}.".format(host))

    SSH_RUN = """
    ssh -o StrictHostKeyChecking=no -o ConnectTimeout=1 {}@{} '
        echo $[`who | cut -d " " -f 1 | sort -u | wc -l`
              + `ps -A | grep main | wc -l`]'
    """

    to_run = SSH_RUN.format(id, host)
    try:
        out = subprocess.check_output(to_run, shell=True, timeout=2)
        val = int(out)
        return val == 1
    except:
        return False

def run_remote(id, host, command, file, server, port):
    print("Run remote.")
    SSH_RUN = """
    where={id}@{host}
    ssh -o "StrictHostKeyChecking no" $where "
      echo 'Disapching remote at $where'

      export GOPATH=~/golang
      cd ~/golang/src/github.com/danalex97/nfsTorrent

      nice nohup python3 remote/job.py -s={server} -p={port} -n={file} {command} > /dev/null 2>&1 &
      # python3 remote/job.py -s={server} -p={port} -n={file} {command}
      echo 'Finished dispaching at $where'
      exit
    "
    """

    to_run = SSH_RUN.format(
        id      = id,
        host    = host,
        command = command,
        file    = file,
        server  = server,
        port    = port
    )
    os.system(to_run)

def process_output(file):
    print(os.path.realpath(file))
    with open(file, 'r') as content_file:
        content = content_file.read()

    lines = content.split("\n")

    ans = {}
    for line in lines:
        for k, v in KEYS.items():
            if v in line and k not in ans:
                ans[k] = float(line.split(":")[1])

    if "red" not in ans:
        return None
    return ans
