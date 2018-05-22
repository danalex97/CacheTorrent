let Node = function(id) {
  function get_domain(id) {
    return id.split(".")[0];
  }

  return {
    id     : id,
    group  : get_domain(id),
    leader : false,
    active : false,
    pieces : 0,
  };
};

let Link = function(src, dst) {
  return {
    source : src,
    target : dst,
    active : false,
  };
};
