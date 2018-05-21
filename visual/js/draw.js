let NodeDrawer = function(simulation, nodes) {
  function draw(ctx) {
    return ctx
      .enter()
      .append("circle")
      .attr("fill", "red")
      .attr("r", 5);
  }

  self = this;

  self.simulation = simulation;

  self.simulation.force = self.simulation.force.nodes(nodes);
  self.nodeGroup = draw(self.simulation.svg
    .append("g")
    .attr("class", "nodes")
    .selectAll("circle")
    .data(nodes));

  self.restart = function(nodes) {
    self.nodeGroup = nodeGroup.data(nodes, function(d) {
      return d.id;
    });
    self.nodeGroup.exit().remove();
    self.nodeGroup = draw(self.nodeGroup).merge(self.nodeGroup);

    self.simulation.force
      .nodes(nodes);
  };

  self.tick = function() {
    self.nodeGroup
      .attr("cx", function(d) {
        return d.x;
      })
      .attr("cy", function(d) {
        return d.y;
      });
  };

  // Add self to simulation drawers.
  self.simulation.drawers.push(self);

  return self;
};

let nodes = [
  Node(),
  Node(),
  Node(),
  Node()
];

let simulation = Simulation();
let nodeDrawer = NodeDrawer(simulation, nodes);

simulation.start();
