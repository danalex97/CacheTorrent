from util import kill_remote
from util import ID

from pool import Pool

import threading
import time

class KillPool(Pool):
    def __init__(self):
        super().__init__()
        self.running = []

    def next(self):
        if len(self.pool) == 0:
            time.sleep(3)
            for thd in self.running:
                thd.join(0)
            return None
        out = self.pool[-1]
        del self.pool[-1]

        thd = threading.Thread(target=kill_remote, args=[ID, out]).start()
        self.running.append(thd)

        return out

def kill_all():
    print("Killing remote jobs...")
    pool = KillPool()
    while pool.next() != None:
        pass

if __name__ == "__main__":
    kill_all()
