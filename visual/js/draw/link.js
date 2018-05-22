let LinkDrawer = function(ctx, links) {
  let self = this;

  function active(ctx) {
    return ctx.attrs(function(d) {
      if (!d.active) {
        return {
          "stroke" : colors.links.inactive,
          "stroke-width" : "1.5px",
          "stroke-opacity" : "0.3"
        }
      } else {
        return {
          "stroke" : colors.links.active,
          "stroke-width" : "2px",
          "stroke-opacity" : "0.8"
        }
      }
    })
  }

  function draw(ctx) {
    return active(ctx.append("line"));
  }

  function restart() {
    // Apply the general update pattern to the links.
    self.link = self.link.data(self.links);
    self.link.exit().remove();
    self.link = draw(self.link.enter()).merge(self.link);

    // Update and restart the simulation.
    ctx.simulation.force("link").links(self.links);
    ctx.simulation.alpha(1).restart();
  }

  /* Fields. */
  self.link = ctx.center
    .append("g")
    .selectAll(".link");
  self.links = links;

  /* Interface. */
  self.tick = function() {
    self.link
      .attr("x1", function(d) { return d.source.x; })
      .attr("y1", function(d) { return d.source.y; })
      .attr("x2", function(d) { return d.target.x; })
      .attr("y2", function(d) { return d.target.y; });

    active(self.link);
  };

  self.addLink = function(link) {
    self.links.push(link);
    restart();
  };

  self.start = function() {
    restart();
  };

  return self;
};
