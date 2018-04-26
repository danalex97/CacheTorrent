import os
import random
import threading
import subprocess

ID  = "ad5915"
EXT = "doc.ic.ac.uk"

class Pool:
    def __init__(self):
        def app(id, idx):
            if idx < 10:
                id = id + "0"
            return id + str(idx)
        POOL = \
            [app("point", i) for i in range(1, 41)] + \
            [app("matrix", i) for i in range(1, 21)] + \
            [app("graphic", i) for i in range(1, 21)]
        self.pool = POOL[:]
        random.shuffle(self.pool)

    def next(self):
        if len(self.pool) == 0:
            return None
        out = self.pool[-1]
        del self.pool[-1]

        if test_remote(ID, out):
            return out
        return self.next()

class Job:
    def __init__(self, pool, command, times):
        self.pool = pool

        self.command = command
        self.times   = times

    def run(self):
        def run(host):
            os.system("mkdir remote_run")
            file = "remote_run/{}.txt".format(host)

            run_remote(ID, host, self.command, file)

            print(process_output(file))

        for _ in range(self.times):
            host = self.pool.next()
            threading.Thread(target=run, args=[host]).start()

def test_remote(id, host):
    SSH_RUN = """
    ssh -t -o StrictHostKeyChecking=no -o ConnectTimeout=1 {}@{} 'who | cut -d " " -f 1 | sort -u | wc -l'
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
    ssh -tt -o "StrictHostKeyChecking no" $where <<-'ENDSSH'
      echo "Running remote at $where"

      export GOPATH=~/golang
      cd ~/golang/src/github.com/danalex97/nfsTorrent

      {} > {}
      exit
    ENDSSH > /dev/null
    """

    to_run = SSH_RUN.format(id, host, command, file)
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

    if not ans["red"]:
        return None
    return ans

if __name__ == "__main__":
    pool = Pool()
    jobs = [
        Job(pool, "go run main.go -ext -conf=confs/small.json", 10),
        Job(pool, "go run main.go -conf=confs/small.json", 10),
    ]
    for job in jobs:
        job.run()
