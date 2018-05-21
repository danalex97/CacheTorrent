Ctx = function() {
  var svg = d3.select("svg"),
      width = +svg.attr("width"),
      height = +svg.attr("height");

  var simulation = d3.forceSimulation()
      .force("charge", d3.forceManyBody())
      .on("tick", ticked);

  return {
    "svg" : svg,
    "width" : width,
    "height" : height,
    "simulation" : simulation
  }
}

function ticked() {
  drawer.tick();
}
