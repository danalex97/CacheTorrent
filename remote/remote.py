#!/usr/bin/python3

import os
import threading
import string
import time
import sys
import argparse

from util import test_remote
from util import random_id

from pool import Pool

ID  = "ad5915"
EXT = "doc.ic.ac.uk"

keys = {
  "red"  : "Redundancy",
  "time" : "Average time",
  "50p"  : "50th percentile",
  "90p"  : "90th percentile",
}

class Job:
    os.system("mkdir -p remote_run")
    os.system("mkdir -p results")

    def __init__(self, pool, command, times, name):
        self.pool = pool

        self.lock = threading.Lock()
        self.results = []

        self.command = command
        self.times   = times
        self.runs    = 0

        self.id = name

        os.system("mkdir -p results/{}".format(self.id))
        os.system("mkdir -p results/{}/runs".format(self.id))

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
                self.lock.acquire()
                self.runs += 1
                runs = self.runs
                self.lock.release()

                out = open("results/{}/runs/{}.txt".format(self.id, runs), 'w')

                print("===========================", file=out)
                print("Job: {} -- run".format(self.command), file=out)
                print("===========================", file=out)
                for k, v in res.items():
                    print("{} : {}".format(keys[k], v), file=out)

                out.close()

            self.lock.acquire()
            self.results.append(res)
            self.lock.release()

        print("Job id: {}".format(self.id))
        print("Running job: {}".format(self.command))
        for _ in range(int(self.times)):
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

    parser.add_argument("-n", "--name", type=str, default=random_id(),
        help="The name of the folder in which the results will be saved.")
    parser.add_argument("-r", "--runs", type=int, default=1,
        help="Number of times that the job runs.")
    parser.add_argument('command', nargs='*')

    args, command_flags = parser.parse_known_args()

    runs    = args.runs
    name    = args.name
    command = " ".join(args.command + command_flags)

    print("Running remote command: {}".format(command))

    # Run the job remotely
    job = Job(Pool(), command, runs, name).run().wait()

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
