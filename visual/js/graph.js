let Node = function(id) {
  function get_domain(id) {
    return id.split(".")[0];
  }

  return {
    id     : id,
    group  : get_domain(id),
  };
};

let Link = function(src, dst) {
  return {
    source : src,
    target : dst,
    active : false,
  };
};
