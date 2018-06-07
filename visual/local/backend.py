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
            print("Job received: {}".format(message))
            if message is None:
                return self.fail(message)
            if "log" not in message:
                return self.fail(message)
            log  = message["log"]
            args = [("log", message["log"])]

            for k, v in message.items():
                if k != "log":
                    args.append((k, v))

            command = "go run main.go {}".format(format_args(args))
            print("> {}".format(command))
            os.system(command)

            return {
                "ok" : "true",
            }

    def __init__(self, port=8080):
        self.server = Server("coordinator", port)
        self.server.add_component_post("/job", Backend.OnJob(self))

    def run(self):
        threading.Thread(target=self.server.run).start()
