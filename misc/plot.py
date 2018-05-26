import matplotlib.pyplot as plt

plt.style.use('classic')
# plt.style.use('seaborn-white')
# plt.style.use('fivethirtyeight')

def setup(plt):
    width, height = plt.figaspect(1.0 / 2.0)
    plt.figure(figsize=(width, height), dpi=400)

    plt.rcParams['figure.facecolor'] = '0.75'
    plt.rcParams['axes.labelsize'] = 15

    plt.rcParams['xtick.labelsize'] = 12
    plt.rcParams['ytick.labelsize'] = 12

    plt.rcParams['legend.fontsize'] = 15

    plt.rcParams['axes.linewidth'] = 1
    plt.rcParams['lines.markeredgewidth'] = 1

    plt.margins(0.02, 0.02)

class Plot():
    def __init__(self):
        self.plt = plt
        setup(self.plt)

        self.loc = 4

    def init(self):
        pass

    def xlabel(self, xlabel):
        self.plt.xlabel(xlabel)
        return self

    def ylabel(self, ylabel):
        self.plt.ylabel(ylabel)
        return self

    def plot(self, *args, **kwargs):
        self.plt.plot(*args, **kwargs)
        return self

    def save(self, name):
        print("Saving {}".format(name))
        self.plt.legend(loc = self.loc)
        self.plt.savefig(name)

class CDFPlot(Plot):
    def __init__(self):
        super().__init__()
        self.xlabel("Full file download time(ms)")
        self.ylabel("Peers with finished download")
        self.loc = 4

class LeaderPlot(Plot):
    def __init__(self):
        super().__init__()
        self.xlabel("Full file download time(ms)")
        self.ylabel("Peer percent with download finished")
        self.plt.xlim((1200000,2000000))
        self.loc = 4

types = {
    "Bias" : (lambda self: self.xlabel("Percent of external connections")),
    "Cache" : (lambda self: self.xlabel("Percent of leaders inside a domain")),
    "Both" : (lambda self: self.xlabel("Percent of leaders inside a domain/external connections")),
}

class TimePlot(Plot):
    def __init__(self, type):
        super().__init__()
        self.ylabel("Full file download time(ms)")
        self.loc = 1
        types[type](self)

class RedPlot(Plot):
    def __init__(self, type):
        super().__init__()
        self.ylabel("Redundant transmissions")
        self.loc = 4
        types[type](self)

    def setLoc(self, loc):
        self.loc = loc
        return self

class PiecePlot(Plot):
    def __init__(self):
        super().__init__()
        self.ylabel("Full file download time(ms)")
        self.xlabel("Number of pieces")
        self.loc = 1
