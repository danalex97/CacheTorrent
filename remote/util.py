import os
import subprocess
import socket

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
}

random_id = lambda : ''.join(
    random.choice(string.ascii_lowercase + string.digits)
    for _ in range(20))

def test_remote(id, host):
    # Check if there is only one user and
    # my script is not already running on the host.
    SSH_RUN = """
    ssh -o StrictHostKeyChecking=no -o ConnectTimeout=1 {}@{} '
        echo $[`who | cut -d " " -f 1 | sort -u | wc -l`
              + `ps -A | grep main | wc -l`]'
    """

    to_run = SSH_RUN.format(id, host)
    try:
        out = subprocess.check_output(to_run, shell=True)
        val = int(out)
        return val == 1
    except:
        return False

def run_remote(id, host, command, file, server, port):
    to_run = """
    where={id}@{host}
    ssh -o "StrictHostKeyChecking no" $where "
      echo 'Running remote at $where'

      export GOPATH=~/golang
      cd ~/golang/src/github.com/danalex97/nfsTorrent

      nohop python3 remote/job.py {command} -s={sender} -p={port} -n={file} > /dev/null 2>&1 &
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
    with open(file, 'r') as content_file:
        content = content_file.read()

    lines = content.split("\n")

    ans = {}
    for line in lines:
        for k, v in keys.items():
            if v in line:
                ans[k] = float(line.split(":")[1])

    if "red" not in ans:
        return None
    return ans
