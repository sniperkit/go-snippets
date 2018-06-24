function createPopup(url, dialogTitle) {
		$.get( url, function(data) {
				var d= $( "<div id='popup'>").html(data).dialog({
						width:"auto", 
						height:"auto", 
						title: dialogTitle,
						position: {
								my:"center bottom",
								at:"center top+50%",
								of:window
						}
				});
		});
}

function toggleVisibility(objId) {
		if ($( objId ).is(":visible")) {
				$( objId).hide();
		} else {
				$( objId).show();
		}
}

function create3DScatterPlot(coordinatesID, labelsID, targetDiv) {
		var lines= $( "#"+coordinatesID ).html().split("\n");
		var labels= $( "#"+labelsID ).html().split("\n");
		data = [];
		for(var i=0;i<lines.length;i++) {
				str = lines[i].split(",");
				var tuple = {};
				var defineTuple = false;
				if (str.length > 0){
						tuple.x = parseFloat(str[0]);
						if (!isNaN(tuple.x)){
								defineTuple = true;
						}
				} else {
						tuple.x = 0.0;
				}
				if (str.length > 1){
						tuple.y = parseFloat(str[1]);
						defineTuple = true;
				} else {
						tuple.y = 0.0;
				}
				if (str.length > 2){
						tuple.z = parseFloat(str[2]);
						defineTuple = true;
				} else {
						tuple.z = 0.0;
				}
				if (defineTuple) {
						tuple.name = labels[i];
						tuple.scoreval= -1.0;
						data.push(tuple);
				}
		}
		var maxElem={x:"",y:"",z:""}, minElem = {x:"", y:"", z:""};
		for(var i=0;i<data.length;i++) {
				o = data[i]
				if (maxElem.x == "" || maxElem.x < o.x) {
						maxElem.x = o.x;
				}
				if (maxElem.y == "" || maxElem.y < o.y) {
						maxElem.y = o.y;
				}
				if (maxElem.z == "" || maxElem.z < o.z) {
						maxElem.z = o.z;
				}
				if (minElem.x == "" || minElem.x > o.x) {
						minElem.x = o.x;
				}
				if (minElem.y == "" || minElem.y > o.y) {
						minElem.y = o.y;
				}
				if (minElem.z == "" || minElem.z > o.z) {
						minElem.z = o.z;
				}

		}
		// Give the points a 3D feel by adding a radial gradient
		Highcharts.getOptions().colors = $.map(Highcharts.getOptions().colors, function (color) {
				return {
						radialGradient: {
								cx: 0.4,
								cy: 0.3,
								r: 0.5
						},
						stops: [
								[0, color],
								[1, Highcharts.Color(color).brighten(-0.2).get('rgb')]
						]
				};
		});

		// Set up the chart
		chart = new Highcharts.Chart({
				chart: {
						renderTo: targetDiv,
						margin: 100,
						type: 'scatter',
						options3d: {
								enabled: true,
								alpha: 10,
								beta: 30,
								depth: 800,
								viewDistance: 5,
								fitToPlot: false,
								frame: {
										bottom: { size: 1, color: 'rgba(0,0,0,0.02)' },
										back: { size: 1, color: 'rgba(0,0,0,0.04)' },
										side: { size: 1, color: 'rgba(0,0,0,0.06)' }
								}
						},
				},
				title: {
						text: 'Dataset space'
				},
				subtitle: {
						text: 'Click and drag the plot area to rotate in space'
				},
				tooltip: {
						pointFormatter: function(){ 
								var message = this.name+"<br/>(";
								message +=parseFloat(this.x)+","
								message +=parseFloat(this.y)+","
								message +=parseFloat(this.z)+")"
								if(scores!=undefined && scores[this.name]!=undefined) {
										message+="<br/>Operator score:"+parseFloat(scores[this.name]);
								}
								return message
						}
				},
				plotOptions: {
						scatter: {
								width: 10,
								height: 10,
								depth: 10
						}
				},
				yAxis: {
						title: { 
								text : "PC2"
						}
				},
				xAxis: {
						title: { 
								text : "PC1"
						},
						gridLineWidth: 1
				},
				zAxis: {
						title: { 
								text : "PC3"
						},
						showFirstLabel: false
				},
				legend: {
						enabled: false
				},
				series: [{
						name: 'Dataset',
						colorByPoint: false,
						data : data
				}]
		});


		// Add mouse events for rotation
		$(chart.container).on('mousedown.hc touchstart.hc', function (eStart) {
				eStart = chart.pointer.normalize(eStart);

				var posX = eStart.pageX,
						posY = eStart.pageY,
						alpha = chart.options.chart.options3d.alpha,
						beta = chart.options.chart.options3d.beta,
						newAlpha,
						newBeta,
						sensitivity = 6; // lower is more sensitive

				$(document).on({
						'mousemove.hc touchdrag.hc': function (e) {
								// Run beta
								newBeta = beta + (posX - e.pageX) / sensitivity;
								chart.options.chart.options3d.beta = newBeta;

								// Run alpha
								newAlpha = alpha + (e.pageY - posY) / sensitivity;
								chart.options.chart.options3d.alpha = newAlpha;

								chart.redraw(false);
						},
						'mouseup touchend': function () {
								$(document).off('.hc');
						}
				});
		});
};

function colorizePoints(obj) {
		colScales = [[2,0,0],[2,1,0],[2,2,0],[0,2,0],[0,2,2],[0,1,2],[0,0,2],[0,0,0]].reverse();
		colorStep = 100;
		id = obj.value;
		if (id == "none") {
				// do nothing
				scores={}
				for(i=0;i<data.length;i++) {
						data[i].color = "" 
				}
				chart.series[0].update({data:data});
				console.log($("#legend").html())
				$("#legend").remove();

				return;
		}
		$.get("/scores/"+id+"/text", function(d){
				scores={};
				lines = d.split("\n");
				minElem = undefined, maxElem =undefined;
				for(i=0;i<lines.length;i++) {
						arr = lines[i].split(":")
						var v =parseFloat(arr[1])
						scores [arr[0]] = v
						if (minElem == undefined || minElem > v) {
								minElem = v
						}
						if (maxElem == undefined || maxElem < v) {
								maxElem = v
						}
				}
				data = chart.series[0].data;
				colorRegions = [];
				for(i=0;i<data.length;i++) {
						s = scores[data[i].name]
						v = (s - minElem)/(maxElem-minElem)
						idx = Math.round(v*(colScales.length-1))
						c = colScales[idx]
						if (colorRegions[idx] == undefined) {
								colorRegions[idx] = {min: s, max:s, count:0}
						}
						if (colorRegions[idx].min > s) {
								colorRegions[idx].min = s
						}
						if (colorRegions[idx].max < s) {
								colorRegions[idx].max = s
						}
						colorRegions[idx].count  =colorRegions[idx].count+1
						if (c!=undefined ){
								r= c[0]*colorStep, g = c[1]*colorStep, b=c[2]*colorStep
								rgbString = "rgb("+parseInt(r)+","+parseInt(g)+","+parseInt(b)+")" 
								data[i].color = rgbString 
						}
				}
				legendDiv = "<table class='tablelist' style='font-weight:bold;'>\n";
				legendDiv = legendDiv+"<tr><th>#</th><th>min</th><th>max</th><th>#points</th></tr>"
				for(var i=0;i<colScales.length;i++) {
						if (c!=undefined ){
								c = colScales[i]
								r= c[0]*colorStep, g = c[1]*colorStep, b=c[2]*colorStep
								rgbString = "rgb("+parseInt(r)+","+parseInt(g)+","+parseInt(b)+")" 
								if (colorRegions[i] == undefined) {
										legendDiv = legendDiv +
												"<tr style='color:"+rgbString+"'>"+
												"<td>"+parseInt(i+1)+"</td>"+
												"<td style='text-align:right;'>-</td>"+
												"<td style='text-align:right;'>-</td>"+
												"<td style='text-align:right;'>0</td>"+
												"</tr> ";

								} else {
										legendDiv = legendDiv +
												"<tr style='color:"+rgbString+"'>"+
												"<td>"+parseInt(i+1)+"</td>"+
												"<td style='text-align:right;'>"+colorRegions[i].min.toFixed(2)+"</td>"+
												"<td style='text-align:right;'>"+colorRegions[i].max.toFixed(2)+"</td>"+
												"<td style='text-align:right;'>"+colorRegions[i].count+"</td>"+
												"</tr> ";
								}
						}
				}
				legendDiv = legendDiv + "</table>";
				$("<div id='legend'></div>").dialog({
						width:"auto", 
						height:"auto", 
						title:"Color Legend",
						position: {
								my:"left top",
								at:"right top",
								of:$("#main"),
						}
				}).html(legendDiv).attr("id", "legend");
				chart.series[0].update({data:data});
		});
}

function createHeatmap(smID) {
		$.get(smID, function(d) {
				lines = d.split("\n")
				var data = []
				idx={}
				labels=[]
				dict = {}
				minVal=1.0;
				for (i=1;i<lines.length;i++) {
						arr = lines[i].split(",")
						if (arr[0]!=undefined && idx[arr[0]] == undefined && arr[0]!="") {
								key = arr[0], ix = Object.keys(idx).length
								idx[key] = ix
								labels[ix] = key
						}
						if ( arr[1]!=undefined && idx[arr[1]] == undefined && arr[1]!="") {
								key = arr[1], ix = Object.keys(idx).length
								idx[key] = ix
								labels[ix] = key
						}
						v = parseFloat(arr[2]);
						if (minVal>v) {
								minVal = v;
						}
						if(idx[arr[0]]!=undefined && idx[arr[1]]!=undefined) {
								if(dict[arr[0]] == undefined) {
										dict[arr[0]] = {}
								}
								dict[arr[0]][arr[1]]=v
								data.push([idx[arr[0]], idx[arr[1]], v])
						}
				}
				chart = Highcharts.chart('container', {
						chart: {
								type: 'heatmap',
								margin: [60, 10, 80, 50]
						},

						boost: {
								useGPUTranslations: true
						},

						title: {
								text: 'Similarity Matrix as a heatmap',
								align: 'left',
								x: 40
						},

						subtitle: {
								text: 'Similarity Matrix',
								align: 'left',
								x: 40
						},
						tooltip: {
								pointFormatter: function() {
										return labels[this.x]+","+labels[this.y]+":"+this.value
								}
						},

						xAxis: {
								labels: {
										formatter: function() {return labels[this.x]}
								}
						},

						yAxis: {
								labels: {
										formatter: function() {return labels[this.y];}
								},
								tickLength: 20
						},

						colorAxis: {
								stops: [
										[0, '#3060cf'],
										[0.5, '#fffbbc'],
										[0.9, '#c4463a'],
										[1, '#c4463a']
								],
								min: minVal,
								max: 1.0,
								startOnTick: true,
								endOnTick: true
						},

						series: [{
								data: data
						}]

				});
		});
};


function changeOrdering(mainDataset) {
		mainDataset=mainDataset.value
		ordered = labels;
		if (mainDataset == "" ) {
				ordered.sort()
		}else {
				ordered.sort(function(a,b){return dict[mainDataset][a] - dict[mainDataset][b]}).reverse()
		}
		newData = []
		for(i=0;i<labels.length;i++) {
				for(j=0;j<labels.length;j++) {
						newData.push([i,j,dict[ordered[i]][ordered[j]]])
				}
		}
		chart.series[0].update({data:newData});
}

function trackDataset(option) {
		var obj;
		var data = chart.series[0].data;
		for(var i=0;i<data.length;i++) {
				obj = data[i];
				obj.marker = null;
				if (data[i].name == option.value) {
						obj.marker = {
								enabled:true,
								radius : 15,
								fillColor: obj.color
						};
				}
		}
		if (obj!=undefined) {
				chart.series[0].update({data:data});
		}
}
