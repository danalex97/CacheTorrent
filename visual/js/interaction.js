function format(res) {
  res = res.replace(/'/g, '"');
  res = JSON.parse(res);
  return res;
}

// Radio button state.
let currentValue = 1;
function handleClick(myRadio) {
  currentValue = myRadio.value;

  $('#option1').removeClass('active');
  $('#option2').removeClass('active');

  $('#option' + currentValue).addClass('active');
}

function amRequesting() {
  return currentValue == 2;
}

// Request state.
function request(data) {
  let where = data.log;

  if (!amRequesting()) {
    get_log(where, main);
    return;
  }

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
