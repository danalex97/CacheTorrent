var svg = d3.select("svg"),
    width = +svg.attr("width"),
    height = +svg.attr("height");

var a = {id: "a"},
    b = {id: "b"},
    c = {id: "c"},
    nodes = [a,b,c];

var simulation = d3.forceSimulation()
    .force("charge", d3.forceManyBody())
    .on("tick", ticked);

function ticked() {
  drawer.tick();
}
