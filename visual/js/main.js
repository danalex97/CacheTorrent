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

  let nodeDrawer = new NodeDrawer(ctx, nodes);
  let linkDrawer = new LinkDrawer(ctx, links);
  ctx.addTicker(linkDrawer);
  ctx.addTicker(nodeDrawer);
  ctx.addStarter(nodeDrawer);
  ctx.addStarter(linkDrawer);

  ctx.start();
}
