import os
import random

ID  = "ad5915"
EXT = "doc.ic.ac.uk"

class Pool:
    def __init__(self):
        def app(id, idx):
            if idx < 10:
                id = id + "0"
            return id + str(idx)
        POOL = \
            [app("point", i) for i in range(0,20)] + \
            [app("matrix", i) for i in range(0,20)] + \
            [app("ray", i) for i in range(0,20)]
        self.pool = POOL[:]
        random.shuffle(self.pool)

    def next(self):
        if len(self.pool) == 0:
            return None
        out = self.pool[-1]
        del self.pool[-1]

        return out


def run_remote(id, host, file):
    SSH_RUN = """
    where={}@{}
    ssh -tt -o "StrictHostKeyChecking no" $where <<-'ENDSSH'
      echo "Running remote at $where"

      export GOPATH=~/golang
      cd ~/golang/src/github.com/danalex97/nfsTorrent

      go run main.go > {}
      exit
    ENDSSH
    """

    to_run = SSH_RUN.format(id, host, file)
    os.system(to_run)

def process_output(file):
    with open(file, 'r') as content_file:
        content = content_file.read()

    lines = content.split("\n")

    keys = {
      "red"  : "Redundancy",
      "time" : "Average time",
      "50p"  : "50th percentile",
      "90p"  : "90th percentile",
    }
    ans = {}
    for line in lines:
        for k, v in keys.items():
            if v in line:
                ans[k] = float(line.split(":")[1])

    print(ans)
    if not ans["red"]:
        return None
    return ans

# for host in HOSTS:
#     file = "hmm.txt"
#     run_remote(ID, host + EXT, file)
