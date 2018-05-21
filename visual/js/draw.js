const width = window.innerWidth;
const height = window.innerHeight;

const svg = d3.select('svg')
  .attr('width', width)
  .attr('height', height);

const simulation = d3.forceSimulation()
  .force('charge', d3.forceManyBody().strength(-20))
  .force('center', d3.forceCenter(width / 2, height / 2))


var nodes_data =  [
  {"name": "Travis", "sex": "M"},
  {"name": "Rake", "sex": "M"},
  {"name": "Diana", "sex": "F"},
  {"name": "Rachel", "sex": "F"},
  {"name": "Shawn", "sex": "M"},
  {"name": "Emerald", "sex": "F"}
]

simulation.nodes(nodes_data)

var node = svg.append("g")
  .attr("class", "nodes")
  .selectAll("circle")
  .data(nodes_data)
  .enter()
  .append("circle")
  .attr("r", 5)
  .attr("fill", "red");

function tickActions() {
  node
    .attr("cx", function(d) { return d.x; })
    .attr("cy", function(d) { return d.y; })
}

simulation.on("tick", tickActions );
