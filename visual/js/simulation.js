Ctx = function() {
  self = this;

  function tick() {
    self.tickers.forEach(function(ticker) {
      ticker.tick();
    });
  }

  self.svg    = d3.select("svg");
  self.width  = +self.svg.attr("width");
  self.height = +self.svg.attr("height");
  self.simulation = d3.forceSimulation()
      .force("charge", d3.forceManyBody())
      .on("tick", tick);

  self.tickers = [];

  self.addTicker = function(ticker) {
    self.tickers.push(ticker);
  }

  return self
}
