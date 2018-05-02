import os
import random
import threading
import subprocess
import string
import time
import sys
import argparse

ID  = "ad5915"
EXT = "doc.ic.ac.uk"

keys = {
  "red"  : "Redundancy",
  "time" : "Average time",
  "50p"  : "50th percentile",
  "90p"  : "90th percentile",
}

random_id = lambda : ''.join(
    random.choice(string.ascii_lowercase + string.digits)
    for _ in range(20))

class Pool:
    def __init__(self):
        def app(id, idx):
            if idx < 10:
                id = id + "0"
            return id + str(idx)
        POOL = \
            [app("point", i) for i in range(1, 41)] + \
            [app("matrix", i) for i in range(1, 41)] + \
            [app("sprite", i) for i in range(1, 41)] + \
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
    os.system("mkdir -p remote_run")
    os.system("mkdir -p results")

    def __init__(self, pool, command, times):
        self.pool = pool

        self.lock = threading.Lock()
        self.results = []

        self.command = command
        self.times   = times
        self.runs    = 0

        self.id = random_id()

        os.system("mkdir results/{}".format(self.id))
        os.system("mkdir results/{}/single".format(self.id))
        os.system("mkdir results/{}/additional".format(self.id))

    def run(self):
        def run(host):
            file = "remote_run/{}.txt".format(host)

            res = None
            try:
                run_remote(ID, host, self.command, file)
                res = process_output(file)
            except:
                pass

            if res != None:
                out = None

                self.lock.acquire()
                self.runs += 1
                runs = self.runs
                self.lock.release()

                if runs <= self.times:
                    out = open("results/{}/single/{}.txt".format(self.id, runs), 'w')
                else:
                    out = open("results/{}/additional/{}.txt".format(self.id, runs), 'w')

                print("===========================", file=out)
                print("Job: {} -- run".format(self.command), file=out)
                print("===========================", file=out)
                for k, v in res.items():
                    print("{} : {}".format(keys[k], v), file=out)

                out.close()

            self.lock.acquire()
            self.results.append(res)
            self.lock.release()

        print("Running job: {}".format(self.command))
        for _ in range(int(self.times * 2.5)):
            host = self.pool.next()
            threading.Thread(target=run, args=[host]).start()
        return self

    def wait(self):
        def get_len():
            self.lock.acquire()
            ln = len(self.results)
            self.lock.release()
            return ln
        while get_len() < self.times:
            time.sleep(1)
        return self

def test_remote(id, host):
    SSH_RUN = """
    ssh -o StrictHostKeyChecking=no -o ConnectTimeout=1 {}@{} 'who | cut -d " " -f 1 | sort -u | wc -l'
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

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Run multiple simulations \
        remotely')

    parser.add_argument("-r", "--runs", type=int, default=1,
        help="Number of times that the job run.")
    parser.add_argument('command', nargs='*')

    args, command_flags = parser.parse_known_args()

    runs    = args.runs
    command = " ".join(args.command + command_flags)

    print("Running remote command: {}".format(command))

    # Run the job remotely
    job = Job(Pool(), command, runs).run().wait()

    # Output the results
    out = open("results/{}/summary.txt".format(job.id), 'w')

    print("===========================", file=out)
    print("Job: {}".format(job.command), file=out)

    rs = list([r for r in job.results if r != None][:job.runs])

    if len(rs) == 0:
        print("Failed!", file=out)
        sys.exit(0)

    print("===========================", file=out)
    print("Summary:", file=out)
    print("===========================", file=out)
    print("Runs: {}".format(job.runs), file=out)
    ans = rs[0]
    for r in rs[1:]:
        for k, v in r.items():
            ans[k] += v
    for k, v in ans.items():
        print("{} : {}".format(keys[k], v / job.runs), file=out)
