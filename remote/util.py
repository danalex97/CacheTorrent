import os
import subprocess

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

def run_remote(id, host, command, file):
    SSH_RUN = """
    where={}@{}
    ssh -o "StrictHostKeyChecking no" $where "
      echo 'Running remote at $where'

      export GOPATH=~/golang
      cd ~/golang/src/github.com/danalex97/nfsTorrent

      {} > {}
      exit
    "
    """

    to_run = SSH_RUN.format(id, host, command, file)
    os.system(to_run)
