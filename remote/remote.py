#!/usr/bin/python3

from util import test_remote
from util import random_id
from util import ID
from util import KEYS as keys

from coordinator import Coordinator

import os
import threading
import string
import time
import sys
import argparse

def onDone(coordinator):
    # Output the results
    out = open("results/{}/summary.txt".format(coordinator.id), 'w')

    print("===========================", file=out)
    print("Job: {}".format(coordinator.command), file=out)

    rs   = coordinator.results
    runs = len(rs)

    if len(rs) == 0:
        print("Failed!", file=out)
        sys.exit(0)

    print("===========================", file=out)
    print("Summary:", file=out)
    print("===========================", file=out)
    print("Runs: {}".format(runs), file=out)
    ans = rs[0]
    for r in rs[1:]:
        for k, v in r.items():
            ans[k] += v
    for k, v in ans.items():
        print("{} : {}".format(keys[k], v / runs), file=out)

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

    # Starting the coordinator
    coordinator = Coordinator(
        command = command,
        times   = runs,
        name    = name) \
    .run()

    coordinator.onDone(lambda: onDone(coordinator))
