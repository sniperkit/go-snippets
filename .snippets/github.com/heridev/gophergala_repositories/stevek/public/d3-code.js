/*
 * Adapted from:
 * http://bl.ocks.org/herrstucki/5684816
 */
var points = [];

var width = 960,
    height = 500;

var pointGrid = d3.layout.grid()
    .points()
    .size([360, 360]);

var svg = d3.select("#viz").append("svg")
    .attr({
        width: width,
        height: height
    })
    .append("g")
    .attr("transform", "translate(70,70)");

function update() {
    var point = svg.selectAll(".point")
        .data(pointGrid(points));
    point.enter().append("circle")
        .attr("class", "point")
        .attr("r", 1e-6)
        .attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"; });
    point.transition()
        .attr("r", 5)
        .attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"; });
    point.exit().transition()
        .attr("r", 1e-6)
        .remove();
}

function push() {
    points.push({});
    update();
}
