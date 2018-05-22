let inputElement = document.getElementById("input");

inputElement.addEventListener("keyup", function(event) {
  event.preventDefault();
  if (event.keyCode === 13) {
    let input = inputElement.value;
    get_log(input, load_log);
  }
});

function get_log(input, callback) {
  $.getJSON(input, function(data) {
    callback(data);
  });
}

function load_log(log) {
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

  let nodes = get_nodes(log);

}
