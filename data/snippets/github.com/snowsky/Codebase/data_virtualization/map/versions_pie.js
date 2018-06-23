var allData = [],
    version = {};

var obj2arr = function (obj) {
    var arr = [];
    for ( var i in obj ) {
        arr.push({ "label": i, "value": obj[i]});
    }
    return arr;
};

var version_graph = function(csvFile, tag, width, height, radius) {
    $("#pie").show();
    $("#machine_list").hide();
    //http://127.0.0.1:60525/map/versions_pie.html?dc=lon&g=pie
    if (window.location.href.indexOf("?") != -1) {
        url = window.location.href.split("?")[1].split("&");
        dc = _.find(url, function(dc) {
            return dc.indexOf("dc") != -1;
        });
        csvFile = "csv/"+dc.split("=")[1]+".csv";
    }

    // Load csv file
    d3.csv(csvFile, function (data) {
        data.forEach(function(line, index) {
            if(line["CenOS Release"] == "") line["CenOS Release"] = "Unknown";
            version[line["CenOS Release"]] = version[line["CenOS Release"]] + 1 || 1;
        });
    //    console.log(version);
        allData = obj2arr(version);

        var pie = new d3pie("pie", {
            header: {
                title: {
                  text: "OS"
                },
                location: "pie-center"
            },
            size: {
                pieInnerRadius: "60%"
            },
            data: {
                sortOrder: "label-asc",
                content: allData
            },
            tooltips:{
                enabled: true,
                type: "placeholder",
                string: "{label}"+"({value})"
            },
            callbacks: {
                onClickSegment: function(a) {
                    $("#pie").hide();
                    $("#machine_list").show();
                    $("#machine_list").append(a.data["label"]+":<br>");
                    data.forEach(function(line, index) {
                        $("#machine_list").append(function() {
                            if(line["CenOS Release"] == a.data["label"])
                                return line["Name"]+"<br>";
                        });
                    });
                }
            }
        });
    });
}

version_graph(csvFile, tag, width, height, radius);