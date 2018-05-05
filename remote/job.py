import argparse
import sys
import os

from util import random_id

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Run a single simulation.')

    parser.add_argument("-n", "--name", type=str, default=None,
        help="The name of the folder in which the results will be saved.")
    parser.add_argument("-s", "--server", type=str, default=None,
        help="The remote server to which the job will send back the results.")
    parser.add_argument('command', nargs='*')

    args, command_flags = parser.parse_known_args()
    print(args)
    print(command_flags)

    if args.name == None:
        print("No name provided.")
        sys.exit(0)
    if args.server == None:
        print("No report server provided. Using 10.0.0.1.")
        args.server = "10.0.0.1"

    name    = args.name
    server  = args.server
    command = args.command

    command_with_flags = " ".join(command + list(command_flags))
    command_with_flags = "{} > {}".format(command_with_flags, name)
    print("Running {}".format(command_with_flags))

    os.system(command_with_flags)
    
