let NodeDrawer = function(ctx, nodes) {
  let self = this;

  function draw(ctx) {
    return ctx
      .append("circle")
      .attr("fill", "red")
      .attr("r", 8);
  }

  function restart() {
    // Apply the general update pattern to the nodes.
    self.node = self.node.data(self.nodes);
    self.node.exit().remove();
    self.node = draw(self.node.enter()).merge(self.node);

    // Update and restart the simulation.
    ctx.simulation.nodes(self.nodes);
    ctx.simulation.alpha(1).restart();
  }

  /* Fields. */
  self.node = ctx.center
    .append("g")
    .attr("stroke", "#fff")
    .attr("stroke-width", 1.5)
    .selectAll(".node");
  self.nodes = nodes;

  /* Interface. */
  self.tick = function() {
    self.node
      .attr("cx", function(d) { return d.x; })
      .attr("cy", function(d) { return d.y; });
  };

  self.addNode = function(node) {
    self.nodes.push(node);
    restart();
  };

  self.start = function() {
    restart();
  };

  return self;
};
