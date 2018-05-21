const width = window.innerWidth;
const height = window.innerHeight;

const svg = d3.select('svg')
  .attr('width', width)
  .attr('height', height);

const simulation = d3.forceSimulation()
  .force('charge', d3.forceManyBody().strength(-20))
  .force('center', d3.forceCenter(width / 2, height / 2));

let nodes = [
  Node(),
  Node(),
  Node(),
  Node(),
  Node()
]

console.log(nodes)
simulation.nodes(nodes);

var node = svg.append("g")
  .attr("class", "nodes")
  .selectAll("circle")
  .data(nodes)
  .enter()
  .append("circle")
  .attr("r", 5)
  .attr("fill", "red");

function tickActions() {
  node
    .attr("cx", function(d) {
      return d.x;
    })
    .attr("cy", function(d) {
      return d.y;
    })
}

d3.timeout(function() {
  nodes.push(Node());
  restart();
}, 1000);

d3.timeout(function() {
  nodes.push(Node());
  restart();
}, 2000);

function restart() {
  // Apply the general update pattern to the nodes.
  node = node.data(nodes, function(d) {
    return d.id;
  });
  node.exit().remove();
  node = node
    .enter()
    .append("circle")
    .attr("fill", "red")
    .attr("r", 5)
    .merge(node);

  // Update and restart the simulation.
  simulation
    .nodes(nodes);
  simulation
    .alpha(1)
    .restart();
}

simulation.on("tick", tickActions);
