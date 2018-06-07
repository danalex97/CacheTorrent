#!/usr/bin/python3

from util import test_remote
from util import random_id
from util import ID
from util import KEYS as keys

from kill import kill_all
from kill import kill_job

from coordinator import Coordinator

import os
import threading
import string
import time
import sys
import argparse

def onDone(coordinator):
    # Output the results
    out = "results/{}/summary.txt".format(coordinator.id)

    # We firstly open the file for writing to overwrite older results
    with open(out, "w") as f:
        print("===========================", file=f)
        print("Job: {}".format(coordinator.command), file=f)

    rs   = coordinator.results
    runs = len(rs)

    if len(rs) == 0:
        with open(out, "a") as f:
            print("Failed!", file=f)
        sys.exit(0)

    with open(out, "a") as f:
        print("===========================", file=f)
        print("Summary:", file=f)
        print("===========================", file=f)
        print("Runs: {}".format(runs), file=f)
    ans = rs[0]
    for r in rs[1:]:
        for k, v in r.items():
            ans[k] += v
    with open(out, "a") as f:
        for k, v in ans.items():
            print("{} : {}".format(keys[k], v / runs), file=f)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Run multiple simulations \
        remotely')

    parser.add_argument("-n", "--name", type=str, default=random_id(),
        help="The name of the folder in which the results will be saved.")
    parser.add_argument("-r", "--runs", type=int, default=0,
        help="Number of times that the job runs.")
    parser.add_argument("-notify", type=int, default=0,
        help="The PID of process to be notified when all jobs were dispached.")
    parser.add_argument("-k", "--kill", nargs='?', action="store", dest="kill", default=[],
        help="Use this flag to kill all remote jobs.")
    parser.add_argument('command', nargs='*')
    parser.set_defaults(kill=False)

    args, command_flags = parser.parse_known_args()

    notify  = args.notify
    runs    = args.runs
    name    = args.name
    kill    = args.kill

    if kill != False:
        if kill == None:
            kill_all()
        else:
            kill_job(str(kill))
        sys.exit(0)

    command = " ".join(args.command + command_flags)

    print("Running remote command: {}".format(command))

    # Starting the coordinator
    coordinator = Coordinator(
        command = command,
        times   = runs,
        name    = name,
        notify  = notify) \

    coordinator.onDone(lambda: onDone(coordinator))
    coordinator.run()
