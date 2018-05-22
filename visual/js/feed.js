let Feed = function (feedArray, interval, nbr) {
  let self = this;

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
    }, interval - interval / 10);
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
  self.interval  = interval;
  self.nbr       = nbr;
  self.pieces    = get_pieces(self.feedArray);

  // Initialize seed.
  self.feedArray[0].source.pieces = self.pieces;

  self.next = function() {
    for (let i = 0; i < nbr; i++) {
      if (self.pos > self.feedArray.length) {
        return;
      }

      process(self.feedArray[self.pos]);
      self.pos += 1;
    }
  };

  self.start = function() {
    d3.interval(self.next, self.interval);
  };

  return self;
};
