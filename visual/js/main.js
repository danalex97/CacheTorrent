var ctx = new Ctx();
var nodes = [new Node(), new Node()];

let drawer  = new NodeDrawer(ctx, nodes);
let drawer2 = new LinkDrawer(ctx, [Link(nodes[0], nodes[1])]);
ctx.addTicker(drawer2);
ctx.addTicker(drawer);
ctx.addStarter(drawer);
ctx.addStarter(drawer2);

ctx.start();
d3.interval(function() {
  drawer.addNode(new Node());

  for (let i = 0; i < 3; i++) {
    let idx1 = Math.floor(Math.random() * nodes.length);
    let idx2 = Math.floor(Math.random() * nodes.length);

    if (idx1 != idx2) {
      drawer2.addLink(Link(nodes[idx1], nodes[idx2]));
    }
  }
}, 1000);
