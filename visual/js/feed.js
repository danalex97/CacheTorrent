let Feed = function (feedArray, interval, nbr) {
  let self = this;

  self.interval  = interval;

  function process(activeLink) {
    if (activeLink.target.pieces == self.pieces) {
      return;
    }

    activeLink.active = true;

    activeLink.target.pieces += 1;
    activeLink.target.active = true;

    d3.timeout(function() {
      activeLink.active = false;
      activeLink.target.active = false;
    }, self.interval - self.interval / 10);
  }

  function get_pieces(feed) {
    let ids    = Array.from(new Set(feed.map(l => l.target.id)));
    let pieces = ids.map(id => feed
      .map(l => l.target.id)
      .filter(l => l == id)
      .length);
    return Math.min.apply(null, pieces);
  }

  self.feedArray = feedArray;
  self.pos       = 0;
  self.nbr       = nbr;
  self.pieces    = get_pieces(self.feedArray);
  self.restart   = false;

  // Initialize seed.
  self.feedArray[0].source.pieces = self.pieces;

  self.next = function() {
    if (self.restart) {
      self.restart = false;
      self.running.stop();

      let interval = self.interval;

      self.running = d3.interval(self.next, interval);
      return;
    }

    for (let i = 0; i < nbr; i++) {
      if (self.pos > self.feedArray.length) {
        return;
      }

      process(self.feedArray[self.pos]);
      self.pos += 1;
    }
  };

  self.start = function() {
    self.running = d3.interval(self.next, self.interval);
  };

  self.setInterval = function(interval) {
    self.interval = interval;
    self.restart = true;
  };

  return self;
};
