/* Default trigger. */
let event = new Event('start');
inputElement.dispatchEvent(event);

/* Main. */
function main(log) {
  d3.select("svg").selectAll("*").remove();

  let env = get_env(log);
  console.log(env);

  let ctx = new Ctx();

  let linkDrawer  = new LinkDrawer(ctx, env.links);
  let nodeDrawer  = new NodeDrawer(ctx, env.nodes);
  let groupDrawer = new GroupDrawer(ctx, env.groups);

  ctx.addTicker(nodeDrawer);
  ctx.addTicker(linkDrawer);
  
  ctx.addStarter(nodeDrawer);
  ctx.addStarter(linkDrawer);
  ctx.addStarter(groupDrawer);

  ctx.start();
}
