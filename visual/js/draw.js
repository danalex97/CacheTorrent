let height = 600;
let width = 1278;

let svg = d3
  .select("body")
  .append("svg")
  .attr("width", width)
  .attr("height", height);

let NodeDrawer = function(node) {
  return node
    .append("circle")
    .attr("class", "node")
    .attr("r", 5)
    .style("fill", "red");
}

let LinkDrawer = function(link) {
  return link
   .append("line")
   .attr("class", "link")
   .style("stroke-width", 5)
}

let GraphDrawer = function() {
  this.force = d3.layout.force()
    .linkDistance(400)
    .size([width, height]);

  this.nodeDrawer = NodeDrawer;
  this.linkDrawer = LinkDrawer;

  this.nodes = [];
  this.links = [];

  this.addNode = function(node) {
    console.log("Add node.")
  }
  this.addLink = function(link) {
    console.log("Add link.")
  }
}
