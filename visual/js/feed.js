let Feed = function (feedArray, interval, nbr) {
  let self = this;

  function process(activeLink) {
    activeLink.active = true;

    activeLink.target.pieces += 1;
    activeLink.target.active = true;

    d3.timeout(function() {
      activeLink.active = false;
      activeLink.target.active = false;
    }, interval - interval / 10);
  }

  self.feedArray = feedArray;
  self.pos       = 0;
  self.interval  = interval;
  self.nbr       = nbr;

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
