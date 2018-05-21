let Simulation = function() {
  self = this;

  self.width = window.innerWidth / 2;
  self.height = window.innerHeight / 2;

  self.svg = d3.select('svg')
    .attr('width', self.width)
    .attr('height', self.height);

  self.force = d3.forceSimulation()
    .force('charge', d3.forceManyBody())
    .force('center', d3.forceCenter(self.width / 2, self.height / 2));

  self.drawers = [];

  self.restart = function() {
    self.drawers.forEach(function(drawer) {
      drawer.restart();
    });
    self.force
      .alpha(1)
      .restart();
  };

  self.start = function() {
    self.force.on("tick", function() {
      self.drawers.forEach(function(drawer) {
        drawer.tick();
      });
    });
  };

  return self;
};
