let GroupDrawer = function(ctx, groups) {
  let self = this;

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

  let coords = getCoords();
  console.log(coords);

  function restart() {
    ctx.simulation.force('x', d3.forceX().x(function(d) {
      return coords[d.group].x;
    }))
    ctx.simulation.force('y', d3.forceY().y(function(d) {
      return coords[d.group].y;
    }))
  }

  /* Interface. */
  self.start = function() {
    restart();
  };

  return self;
};
