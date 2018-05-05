from pool import Pool
from util import get_ip
from util import run_remote
from util import process_output

from server.component import Component
from server.server    import Server

import threading
import os

class OnDone(Component):
    """
    A component that reacts to done/fail notifications.
    """
    def __init__(self, coordinator):
        self.coordinator = coordinator

    def process(self, message):
        if "fail" in message:
            print("Fail: {}".format(message["fail"]))
            self.coordinator.fail(message["fail"])
        elif "done" in message:
            print("Done: {}".format(message["done"]))
            self.coordinator.done(message["done"])
        else:
            # Unexpected message.
            pass

class Coordinator:
    """
    Remote job coordinator. It manages individual jobs.
    """

    # Create the location for remote_run are results
    os.system("mkdir -p remote_run")
    os.system("mkdir -p results")

    def __init__(self, command, times, name, pool=Pool()):
        self.pool = pool

        self.lock = threading.Lock()
        self.results = []

        self.command = command
        self.times   = times

        self.runs      = 0
        self.completed = 0
        self.callback  = lambda *args: None

        self.id = name

        # Prepare file system
        os.system("mkdir -p results/{}".format(self.id))
        os.system("mkdir -p results/{}/runs".format(self.id))

        # Run server
        self.server = get_ip()
        self.port   = 8080

        server = Server("coordinator", self.port)
        server.add_component_post("/done", OnDone(self))
        threading.Thread(target=server.run).start()

    def run(self):
        def run(host):
            """
            Runs a Job at a particular host.
            """
            file = "remote_run/{}.txt".format(host)

            try:
                run_remote(
                    id      = ID,
                    host    = host,
                    command = self.command,
                    file    = file,
                    server  = self.server,
                    port    = self.port,
                )
                self.runs += 1
            except:
                pass

        print("Job id: {}".format(self.id))
        print("Running job: {}".format(self.command))
        for _ in range(int(self.times)):
            host = self.pool.next()
            if host == None:
                continue
            run(host)
            # threading.Thread(target=run, args=[host]).start()
        return self

    def fail(self, file):
        self.completed += 1
        if self.completed == self.runs:
            self.callback()

    def done(self, file):
        self.results.append(process_output(file))
        self.fail()

    def onDone(self, callback):
        self.callback = callback
