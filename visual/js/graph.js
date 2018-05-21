let Node = function(id, x, y) {
  this.id = id;
  this.x = x;
  this.y = y;
}

let Link = function(n1, n2) {
  this.src = n1;
  this.dst = n2;
}
