#!/usr/bin/python3

from util import random_id

from server.client import Client

import argparse
import sys
import os

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Run a single simulation.')

    parser.add_argument("-n", "--name", type=str, default=None,
        help="The name of the folder in which the results will be saved.")
    parser.add_argument("-s", "--server", type=str, default=None,
        help="The remote server to which the job will send back the results.")
    parser.add_argument("-p", "--port", type=int, default=8080,
        help="The remote port to which the job will send back the results.")
    parser.add_argument('command', nargs='*')

    print("Job started.")
    args, command_flags = parser.parse_known_args()

    if args.name == None:
        print("No name provided.")
        sys.exit(0)
    if args.server == None:
        print("No report server provided. Using 10.0.0.1.")
        args.server = "10.0.0.1"

    name    = args.name
    server  = args.server
    port    = args.port
    command = args.command

    # Client used for server communication
    client = Client(server, port)
    client.post("/start", {
        "done" : name
    })

    command_with_flags = " ".join(command + list(command_flags))
    # command_with_flags = "{} > {}".format(command_with_flags, name)
    print("Running {}".format(command_with_flags))

    # Run the command
    exit = os.system(command_with_flags) >> 8

    # Send the response back to the main server
    if exit == 0:
        print("Job done.")
        client.post("/done", {
            "done" : name
        })
    else:
        print("Job failed.")
        client.post("/done", {
            "fail" : name
        })
