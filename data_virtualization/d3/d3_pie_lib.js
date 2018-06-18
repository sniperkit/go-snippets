var allData = [],
    version = {};

var width = height = 1500,
    radius = 400,
    csvFile = "csv/cam.csv",
    tag = "body";

var obj2arr = function (obj) {
    var arr = [];
    for ( var i in obj ) {
        arr.push({ "label": i+"("+obj[i]+")", "value": obj[i]});
    }
    return arr;
};

// Load csv file
d3.csv(csvFile, function (data) {
    data.forEach(function(line, index) {
        if(line["CenOS Release"] == "") line["CenOS Release"] = "Unknown";
        version[line["CenOS Release"]] = version[line["CenOS Release"]] + 1 || 1;
    });
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
        string: "{label}"
      }
    });
});