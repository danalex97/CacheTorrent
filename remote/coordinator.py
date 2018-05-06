from pool import Pool

from util import ID
from util import KEYS as keys

from util import get_ip
from util import run_remote
from util import process_output

from server.component import Component
from server.server    import Server

import threading
import random
import signal
import os

class OnDone(Component):
    """
    A component that reacts to done/fail notifications.
    """
    def __init__(self, coordinator):
        self.coordinator = coordinator

    def process(self, message):
        if "fail" in message:
            with open(self.coordinator.log, "a") as f:
                print("Fail: {}".format(message["fail"]), file=f)
            self.coordinator.fail(message["fail"])
        elif "done" in message:
            with open(self.coordinator.log, "a") as f:
                print("Done: {}".format(message["done"]), file=f)
            self.coordinator.done(message["done"])
        else:
            # Unexpected message.
            pass

class OnStart(Component):
    """
    A component that reacts to start notifications.
    """
    def __init__(self, coordinator):
        self.coordinator = coordinator

    def process(self, message):
        if "start" in message:
            with open(self.coordinator.log, "a") as f:
                print("Start: {}".format(message["start"]), file=f)

class Coordinator:
    """
    Remote job coordinator. It manages individual jobs.
    """

    # Create the location for remote_run are results
    os.system("mkdir -p remote_run")
    os.system("mkdir -p results")

    def __init__(self, command, times, name, notify, pool=Pool()):
        self.pool = pool

        self.lock = threading.Lock()
        self.results = []

        self.command = command
        self.times   = times

        # Notifications
        self.notify = notify

        # Job control
        self.dispaching = True
        self.runs       = 0
        self.completed  = 0
        self.lock       = threading.Lock()

        self.callback  = lambda *args: None

        self.id = name

        # Prepare file system
        os.system("mkdir -p results/{}".format(self.id))
        os.system("mkdir -p results/{}/runs".format(self.id))

        # Logging
        self.log = "results/{}/log.txt".format(self.id)
        os.system("rm -f {}".format(self.log))

        # Run server
        self.ip   = get_ip()
        self.port = random.uniform(30000, 30100)

        self.server = Server("coordinator", self.port)
        self.server.add_component_post("/done", OnDone(self))
        self.server.add_component_post("/start", OnStart(self))
        threading.Thread(target=self.server.run).start()

    def run(self):
        def run(host):
            """
            Runs a Job at a particular host.
            """
            file = "remote_run/{}.txt".format(host)
            with open(self.log, "a") as f:
                print("Run remote job on: {}".format(host), file=f)

            try:
                run_remote(
                    id      = ID,
                    host    = host,
                    command = self.command,
                    file    = file,
                    server  = self.ip,
                    port    = self.port,
                )
                with self.lock:
                    self.runs += 1
            except Exception as e:
                with open(self.log, "a") as f:
                    print("Job failed on: {}".format(host), file=f)
                    print("Exception: {}".format(e), file=f)

        with open(self.log, "a") as f:
            print("Job id: {}".format(self.id), file=f)
            print("Running job: {}".format(self.command), file=f)
        for _ in range(int(self.times)):
            host = self.pool.next()
            if host == None:
                continue
            run(host)
            # threading.Thread(target=run, args=[host]).start()

        # Notify that everything was dispached
        if self.notify != 0:
            os.kill(self.notify, signal.SIGUSR1)

        with self.lock:
            self.dispaching = False
        self.check_finished()

        return self

    def check_finished(self):
        with self.lock:
            if self.completed == self.runs and not self.dispaching:
                # Run callback.
                with open(self.log, "a") as f:
                    print("Jobs finished. Starting callback.", file=f)
                self.callback()

                # Stop server.
                with open(self.log, "a") as f:
                    print("Stopping server.", file=f)
                self.server.stop()

    def job_stats(self, file, res):
        with self.lock:
            runs = self.completed

        out = open("results/{}/runs/{}.txt".format(self.id, runs), 'w')

        print("===========================", file=out)
        print("Job: {} -- run".format(self.command), file=out)
        print("Machine: {}".format(file), file=out)
        print("===========================", file=out)
        for k, v in res.items():
            print("{} : {}".format(keys[k], v), file=out)

        out.close()

    def fail(self, file):
        with self.lock:
            self.completed += 1
        self.check_finished()

    def done(self, file):
        res = process_output(file)
        self.results.append(res)

        self.job_stats(file, res)
        self.fail(file)

    def onDone(self, callback):
        self.callback = callback
