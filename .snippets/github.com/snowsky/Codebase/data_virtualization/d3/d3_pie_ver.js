var width = height = 1000,
    radius = 300,
    csvFile = "csv/cam.csv",
    tag = "body";

drawPie(csvFile, tag, width, height, radius);

var obj2arr = function (obj) {
    var arr = [];
    for ( var i in obj ) {
        arr.push({ "Version": i, "Number": obj[i]});
    }
    return arr;
};

function drawPie ( csvFile, tag, width, height, radius ) {
    var allData = [],
        version = {};
    // Load csv file
    d3.csv(csvFile, function (data) {
        data.forEach(function(line, index) {
            if(line["CenOS Release"] == "") line["CenOS Release"] = "Unknown";
            version[line["CenOS Release"]] = version[line["CenOS Release"]] + 1 || 1;
        });
        allData = obj2arr(version);
        console.log(allData);

        var d3color = d3.scale.category20c();

        // 选择tag
        var svg = d3.select(tag)
            .append("svg")
            .data(allData)
            .attr("width", width)
            .attr("height", height)
            .append("g")
            .attr("transform", "translate(" + radius + "," + radius + ")");

        var arc = d3.svg.arc()
            .outerRadius(radius);

        var pie = d3.layout.pie()
            .value(function(d) {
                return d["Number"];
            });

        var arcs = svg.selectAll("g.slice")
            .data(pie(allData))
            .enter()
            .append("g")
            .attr("class", "slice");    //allow us to style things in the slices (like text)

        arcs.append("path")
            .attr("fill", function(d, i) { return d3color(i); } )
            .attr("d", arc);

        arcs.append("text")                                     //add a label to each slice
            .attr("transform", function(d) {                    //set the label's origin to the center of the arc
                //we have to make sure to set these before calling arc.centroid
                d.innerRadius = 100;
                d.outerRadius = radius;
                return "translate(" + arc.centroid(d) + ")";        //this gives us a pair of coordinates like [50, 50]
            })
            .attr("text-anchor", "middle")                          //center the text on it's origin
            .text(function(d, i) { console.log(allData[i]);return allData[i]["Version"] + "(" + allData[i]["Number"] + ")"; });        //get the label from our original data array
    });
}