let NodeDrawer = function(ctx, nodes) {
  let self = this;

  function draw(toDraw) {
    function dragstarted(d) {
      if (!d3.event.active) {
        ctx.simulation.alphaTarget(0.3).restart();
      }
      d.fx = d.x;
      d.fy = d.y;
    }

    function dragged(d) {
      d.fx = d3.event.x;
      d.fy = d3.event.y;
    }

    function dragended(d) {
      if (!d3.event.active) {
        ctx.simulation.alphaTarget(0);
      }
      d.fx = null;
      d.fy = null;
    }

    function leader(ctx) {
      return ctx.attrs(function(d) {
        if (!d.leader) {
          return {
            "stroke" : "#9ecae1",
            "stroke-width" : "1px",
            "stroke-opacity" : "1",
            "fill" : "#3182bd",
          }
        } else {
          return {
            "stroke" : "#9ecae1",
            "stroke-width" : "1px",
            "stroke-opacity" : "1",
            "fill" : "#ff8c00"
          }
        }
      })
    }

    // Make a group with circle and text.
    let group = toDraw
        .append("g")
        .attr("class", "node");
    // Add text to the group.
    let circle = leader(group.append("circle"))
      .attr("r", 20)
      .call(d3.drag()
          .on("start", dragstarted)
          .on("drag", dragged)
          .on("end", dragended));
    // Add circle to the group.
    let text = group
      .append("text")
      .attr("text-anchor", "middle")
      .attr("dx", 12)
      .attr("dy", ".35em")
      .attr("color", "black")
      .text(function(d) { return d.id });
    return group;
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
    .selectAll(".node");
  self.nodes = nodes;

  /* Interface. */
  self.tick = function() {
    let getX = function(d) {
      if (d.x < -ctx.width * 0.4) { d.x = -ctx.width * 0.4; }
      if (d.x >  ctx.width * 0.4) { d.x =  ctx.width * 0.4; }
      return d.x;
    };
    let getY = function(d) {
      if (d.y < -ctx.height * 0.4) { d.y = -ctx.height * 0.4; }
      if (d.y >  ctx.height * 0.4) { d.y =  ctx.height * 0.4; }
      return d.y;
    };

    self.node
      .selectAll("circle")
      .attr("cx", getX)
      .attr("cy", getY);
    self.node
      .selectAll("text")
      .attr("x", getX)
      .attr("y", getY);
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
