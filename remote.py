import os
import random
import threading
import subprocess
import time

ID  = "ad5915"
EXT = "doc.ic.ac.uk"

keys = {
  "red"  : "Redundancy",
  "time" : "Average time",
  "50p"  : "50th percentile",
  "90p"  : "90th percentile",
}

class Pool:
    def __init__(self):
        def app(id, idx):
            if idx < 10:
                id = id + "0"
            return id + str(idx)
        POOL = \
            [app("point", i) for i in range(1, 41)] + \
            [app("matrix", i) for i in range(1, 41)] + \
            [app("graphic", i) for i in range(1, 41)] + \
            [app("voxel", i) for i in range(1, 41)] + \
            [app("edge", i) for i in range(1, 41)]
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

        self.lock = threading.Lock()
        self.results = []

        self.command = command
        self.times   = times

    def run(self):
        def run(host):
            os.system("mkdir remote_run")
            file = "remote_run/{}.txt".format(host)

            run_remote(ID, host, self.command, file)

            res = process_output(file)

            self.lock.acquire()
            self.results.append(res)
            self.lock.release()

        for _ in range(int(self.times * 1.5)):
            host = self.pool.next()
            threading.Thread(target=run, args=[host]).start()

    def wait(self):
        def get_len():
            self.lock.acquire()
            ln = len(self.results)
            self.lock.release()
            return ln
        while get_len() < self.times:
            time.sleep(1)

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
        Job(pool, "go run main.go -conf=confs/small.json", 10),
        Job(pool, "go run main.go -ext -conf=confs/small.json", 10),
    ]
    for job in jobs:
        job.run()
    for job in jobs:
        job.wait()
    time.sleep(5)
    print("\n\n\n\n\n")
    for job in jobs:
        print("\n")
        print("===========================")
        print("Job: {}".format(job.command))

        rs = list([r for r in job.results if r != None][:job.times])

        if len(rs) < job.times:
            print("Failed!")
            continue

        # print("===========================")
        # print("Runs:")
        # print("===========================")
        # idx = 0
        # for r in rs:
        #     print("Run {}".format(idx))
        #     idx += 1
        #
        #     for k, v in r.items():
        #         print("{} : {}".format(keys[k], v))

        print("===========================")
        print("Summary:")
        print("===========================")
        ans = rs[0]
        for r in rs[1:]:
            for k, v in r.items():
                ans[k] += v
        for k, v in ans.items():
            print("{} : {}".format(keys[k], v / job.times))
