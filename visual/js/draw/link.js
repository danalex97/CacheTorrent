let LinkDrawer = function(ctx, links) {
  let self = this;

  let center = ctx.svg
    .append("g")
    .attr("transform", "translate(" +
      ctx.width / 2 + "," +
      ctx.height / 2 + ")");

  function draw(ctx) {
    return ctx
      .append("line");
  }

  function restart() {
    // Apply the general update pattern to the links.
    self.link = self.link.data(self.links);
    self.link.exit().remove();
    self.link = draw(self.link.enter()).merge(self.link);

    // Update and restart the simulation.
    ctx.simulation.nodes(self.nodes);
    ctx.simulation.alpha(1).restart();
  }

  /* Fields. */
  self.link = center
    .append("g")
    .selectAll(".link");
  self.links = links;

  /* Interface. */
  self.tick = function() {
    self.link
      .attr("x1", function(d) { return d.src.x; })
      .attr("y1", function(d) { return d.src.y; })
      .attr("x2", function(d) { return d.dst.x; })
      .attr("y2", function(d) { return d.dst.y; });
  };

  self.addLink = function(link) {
    self.nodes.push(link);
    restart();
  };

  self.start = function() {
    restart();
  };

  return self;
};
