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
  let curve = d3.curveCardinalClosed.tension(0.3);
  let drawLine = d3.line().curve(curve);

  function draw(toDraw) {
     return toDraw
       .attr("fill-opacity", colors.groupOpacity)
       .attr("fill", colors.group)
       .attr("stroke-opacity", 1)
       .classed('hull', true);
  }

  function get_hull(id) {
    let bound = 50;

    let coords = nodeDrawer.node
      .filter(d => d.group == id)
      .data()
      .map(d => [
        [d.x - bound, d.y],
        [d.x, d.y - bound],
        [d.x + bound, d.y],
        [d.x, d.y + bound]
      ])
      .reduce((arr, curr) => arr.concat(curr), []);

    return d3.polygonHull(coords);
  }

  function restart() {
    // Create force centers.
    ctx.simulation.force('x', d3.forceX().x(function(d) {
      return coords[d.group].x;
    }));
    ctx.simulation.force('y', d3.forceY().y(function(d) {
      return coords[d.group].y;
    }));

    // Calculate convex hulls.
    self.hulls = groups.map(id => get_hull(id));
  }


  /* Fields. */
  // We make a group of other elements.
  self.groups = ctx.center
    .append('g')
    .attr('class', 'groups');
  // We make a hull for each group.
  // Since the number of groups remains unchanged, we make this static.
  self.hull = self.groups
    .selectAll(".hulls")
    .data(groups)
    .enter()
    .append('g')
    .attr('class', 'hulls')
    .append('path');

  /* Interface. */
  self.start = function() {
    restart();
  };

  self.tick = function() {
    // Update data for each id.
    groups.forEach(function(groupId) {
      draw(self.hull
        .filter(id => id == groupId))
        .attr('d', function(id) {
          let points = get_hull(id);
          return drawLine(points);
        });
    });
  };

  return self;
};
