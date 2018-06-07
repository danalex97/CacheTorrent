function format(res) {
  res = res.replace(/'/g, '"');
  res = JSON.parse(res);
  return res;
}

function request(data) {
  let where = data.log;
  $.ajax({
    method: "POST",
    url: "http://127.0.0.1:8080/job",
    data: JSON.stringify(data),
    contentType: 'application/json;charset=UTF-8'
  })
  .done(function(msg) {
    msg = format(msg);
    console.log("Response: ", msg);
    if (msg.ok === 'true') {
      get_log(where, main);
    }
  });
}

d3.select("#bit")
  .on("click", function() {
    request({
      log : "logs/torrent.json"
    });
  });


d3.select("#cache")
  .on("click", function() {
    request({
      log : "logs/cache.json"
    });
  });

d3.select("#bias")
  .on("click", function() {
    request({
      log : "logs/bias.json"
    });
  });
