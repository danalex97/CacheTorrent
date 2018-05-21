Drawer = function(ctx, nodes) {
  let self = this;

  let center = ctx.svg
    .append("g")
    .attr("transform", "translate(" +
      ctx.width / 2 + "," +
      ctx.height / 2 + ")");
  self.node  = center
    .append("g")
    .selectAll(".node");
  self.nodes = nodes;

  self.tick = function() {
    self.node
      .attr("cx", function(d) { return d.x; })
      .attr("cy", function(d) { return d.y; });
  }

  self.restart = function() {
    // Apply the general update pattern to the nodes.
    self.node = node.data(self.nodes);
    self.node.exit().remove();
    self.node = self.node.enter().append("circle").attr("fill", "red").attr("r", 8).merge(self.node);

    // Update and restart the simulation.
    ctx.simulation.nodes(self.nodes);
    ctx.simulation.alpha(1).restart();
  }

  return self;
}
