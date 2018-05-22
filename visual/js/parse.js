function trigger() {
  let input = inputElement.value;
  get_log(input, main);
}

let inputElement = document.getElementById("input");

/* Event listeners. */
inputElement.addEventListener("keyup", function(event) {
  event.preventDefault();
  if (event.keyCode === 13) {
    trigger();
  }
});

inputElement.addEventListener("start", function(event) {
  trigger();
});

/* Functions. */
function get_log(input, callback) {
  $.getJSON(input, function(data) {
    callback(data);
  });
}

function get_env(log) {
  function valid(value) {
    return value != "<invalid Value>";
  }

  function get_nodes(log) {
    let nodes = new Set();
    log.forEach(function(entry) {
      if (valid(entry.src)) {
        nodes.add(entry.src);
      }
      if (valid(entry.dst)) {
        nodes.add(entry.dst);
      }
    });
    return Array.from(nodes).map(id => new Node(id));
  }

  function get_links(log, node_map) {
    let links = new Set();
    log.forEach(function(entry) {
      if (valid(entry.src) && valid(entry.dst)) {
        let src = entry.src;
        let dst = entry.dst;
        links.add(src + ":" + dst);
      }
    });
    return Array.from(links).map(function (link) {
      let src = node_map[link.split(":")[0]];
      let dst = node_map[link.split(":")[1]];
      return new Link(src, dst);
    });
  }

  let nodes = get_nodes(log);
  let links = get_links(log, nodes.reduce(function(map, obj) {
    map[obj.id] = obj;
    return map;
  }, {}));

  return {
    "nodes" : nodes,
    "links" : links,
  }
}
