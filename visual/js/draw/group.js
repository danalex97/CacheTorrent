let GroupDrawer = function(ctx, groups, nodeDrawer) {
  let self = this;

  /* Centers on X-axis. */
  let width  = ctx.width;
  let height = ctx.height;

  function getCoords() {
    let centers = {};
    let n = groups.length;
    for (let i = 0; i < n; i++) {
      let group = groups[i];
      centers[group] = {
        x : (width / (n - 1)) * i - width / 2,
        y : 0,
      };
    }
    return centers;
  }

  // Note: the groupDrawer doesn't support more groups.
  let coords = getCoords();

  /* Curve. */
  function get_hull(id) {
    let coords = nodeDrawer.node
      .filter(d => d.group == id)
      .data()
      .map(d => [d.x, d.y]);

    return d3.polygonHull(coords);
  }

  let curve = d3.curveCardinalClosed.tension(0.3);
  let drawLine = d3.line().curve(curve);

  function draw(toDraw) {
     return toDraw
       .append('path')
       .attr("fill-opacity", 0.1)
       .attr("fill", "red")
       .attr("stroke-opacity", 1)
       .classed('hull', true)
       .attr('d', function(points){
         return drawLine(points);
       });
  }

  function restart() {
    // Create force centers.
    ctx.simulation.force('x', d3.forceX().x(function(d) {
      return coords[d.group].x;
    }))
    ctx.simulation.force('y', d3.forceY().y(function(d) {
      return coords[d.group].y;
    }))

    // Calculate convex hulls.
    self.hulls = groups.map(id => get_hull(id));

    // Draw convex hulls.
    self.hull = self.hull.data(self.hulls);
    self.hull.exit().remove();
    self.hull = draw(self.hull.enter()).merge(self.hull);
  }

  /* Fields. */
  self.hull = ctx.center
    .append("g")
    .selectAll('.hulls');

  /* Interface. */
  self.start = function() {
    restart();
  };

  self.tick = function() {
  };

  return self;
};
