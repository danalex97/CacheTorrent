from util import test_remote
from util import ID

import random

class Pool:
    def __init__(self):
        def app(id, idx):
            if idx < 10:
                id = id + "0"
            return id + str(idx)
        POOL = \
            [app("matrix", i) for i in range(1, 45)] + \
            [app("sprite", i) for i in range(1, 39)] + \
            [app("arc", i) for i in range(1, 15)] + \
            [app("line", i) for i in range(1, 28)] + \
            [app("edge", i) for i in range(1, 41)] + \
            [app("point", i) for i in range(1, 61)] + \
            [app("voxel", i) for i in range(1, 27)] + \
            [app("graphic", i) for i in range(1, 13)]
        self.pool = POOL[:]
        random.shuffle(self.pool)

    def next(self):
        if len(self.pool) == 0:
            return None
        out = self.pool[-1]
        del self.pool[-1]

        if test_remote(ID, out):
            return out
        return self.next()
