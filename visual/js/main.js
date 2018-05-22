/* Default trigger. */
let event = new Event('start');
inputElement.dispatchEvent(event);

/* Main. */
function main(log) {
  d3.select("svg").selectAll("*").remove();

  let env = get_env(log);
  console.log(env);

  let ctx = new Ctx();
  let nodes = env.nodes;
  let links = env.links;

  let linkDrawer = new LinkDrawer(ctx, links);
  let nodeDrawer = new NodeDrawer(ctx, nodes);
  ctx.addTicker(nodeDrawer);
  ctx.addTicker(linkDrawer);
  ctx.addStarter(nodeDrawer);
  ctx.addStarter(linkDrawer);

  ctx.start();
}
