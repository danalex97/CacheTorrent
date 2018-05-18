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
  let self = this;

  self.force = d3.layout.force()
    .linkDistance(400)
    .size([width, height]);

  self.nodeDrawer = NodeDrawer;
  self.linkDrawer = LinkDrawer;

  self.nodes = [];
  self.links = [];

  self.force
    .nodes(self.nodes)
    .links(self.links)
    .start();

  self.addNode = function(node) {
    self.nodes.push(node);
    console.log("Add node.");

    let data = svg
      .selectAll(".node")
      .data(self.nodes)
      .enter()
    self.nodeDrawer(data)
        .call(force.drag);
    return self;
  }

  self.addLink = function(link) {
    self.links.push(link);
    console.log("Add link.");

    let data = svg
      .selectAll(".link")
      .data(self.links)
      .enter()
    self.linkDrawer(data);

    return self;
  }

  return self;
}
