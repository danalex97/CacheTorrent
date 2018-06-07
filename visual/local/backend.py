from remote.server.component import Component
from remote.server.server    import Server

import threading
import random
import signal
import os

def format_args(args):
    sargs = ""
    for k, v in args:
        sargs += " -{}={}".format(k, v)
    return sargs

class Backend:
    """
    Sever that handles front-end requests.
    """
    class OnJob(Component):
        """
        A component that reacts to queries.
        """
        def __init__(self, backend):
            self.backend = backend

        def fail(self, message):
            print("Request {} failed.".format(message))
            return {
                "ok" : "false",
            }

        def process(self, message):
            if "log" not in message:
                fail(message)
            log  = message["log"]
            args = []

            for k, v in message:
                if k != "log":
                    args.append((k, v))
            os.system("go run main.go {}".format(format_args(args)))

            return {
                "ok" : "true",
            }

    def __init__(self, port=8080):
        self.server = Server("coordinator", port)
        self.server.add_component_post("/job", Backend.OnJob(self))

    def run(self):
        threading.Thread(target=self.server.run).start()
