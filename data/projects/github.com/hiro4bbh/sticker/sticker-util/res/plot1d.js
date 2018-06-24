Plot1d = function(container) {
  this.container = container;
  this.svg = d3.select(this.container).html('').append('svg');
  this._marginTop = 0;
  this._marginRight = 70;
  this._marginBottom = 0;
  this._marginLeft = 70;
  var containerPaddingLeft = window.getComputedStyle(this.container, null).getPropertyValue('padding-left');
  var containerPaddingRight = window.getComputedStyle(this.container, null).getPropertyValue('padding-right');
  containerPaddingLeft = +(containerPaddingLeft.substr(0, containerPaddingLeft.length - 2));
  containerPaddingRight = +(containerPaddingRight.substr(0, containerPaddingRight.length - 2));
  this._width = +this.container.clientWidth - (containerPaddingLeft + containerPaddingRight);
  this._height = (+this._width)*1/2;
  this._colorName = '';
  this._colorOnName = '';
  this._nameName = '';
  this._xAxisScale = 'Linear';
};
Plot1d.prototype.color = function(colorName) {
  this._colorName = colorName;
  return this;
};
Plot1d.prototype.colorOn = function(colorOnName) {
  this._colorOnName = colorOnName;
  return this;
};
Plot1d.prototype.data = function(data) {
  this._data = data;
  return this;
};
Plot1d.prototype.name = function(nameName) {
  this._nameName = nameName;
  return this;
};
Plot1d.prototype.xAxisScale = function(scale) {
  this._scale = scale;
  return this;
};
Plot1d.prototype.draw = function(xName) {
  var self = this;
  this._data.sort(function(d, d_) { return d3.ascending(d[xName], d_[xName]) });
  var svg = this.svg;
  svg.attr('width', this._width)
     .attr('height', this._height);
  var svg = svg.append('g').attr('transform', 'translate(' + this._marginLeft + ',' + this._marginTop + ')');
  var tooltip = d3.select(this.container).append('div')
                                         .attr('class', 'tooltip')
                                         .style('opacity', 0);
  tooltip.append('div').attr('class', 'arrow');
  var tooltipInner = tooltip.append('div').attr('class', 'tooltip-inner');
  var width = this._width - (this._marginLeft + this._marginRight);
  var height = this._height - (this._marginTop + this._marginBottom);
  var calculateColorOff = function(d, i) {
    return self._colorName == '' ? '' : d[self._colorName];
  }, calculateColorOn = function(d, i) {
    return self._colorOnName == '' ? calculateColorOff(d, i) : d[self._colorOnName];
  }, returnText = function(d) {
    return self._nameName == '' ? x(d[xName]) : d[self._nameName];
  };
  var pointCircleOn;
  var generateMouseHandler = function(on) {
    return function(d, i) {
      var color = on ? calculateColorOn(d, i) : calculateColorOff(d, i);
      d3.select(this.parentNode).select('#pointLine'+i).attr('stroke-width', on ? 3 : 1);
      var y = d3.select(this.parentNode).select('#pointText'+i).attr('fill', color).attr('y');
      pointCircleOn.attr('cx', x(d[xName]))
                   .attr('fill', color)
                   .style('display', on ? 'block' : 'none');
      tooltip.transition().duration(on ? 250 : 1000).style("opacity", on ? 1.0 : 0.0);
      tooltipInner.text(d[xName]);
      var inTop = y < height/2;
      tooltip.classed('bs-tooltip-top', !inTop)
             .classed('bs-tooltip-bottom', inTop)
             .style('left', (self._marginLeft + x(d[xName]) + 13 - tooltip.node().getBoundingClientRect().width/2) + "px")
             .style('top', (height/2 + (inTop ? 8 : -(tooltip.node().getBoundingClientRect().height + 8))) + "px");
    }
  };
  // Set X-axis
  var x = d3['scale'+this._xAxisScale]().range([0, width]);
  x.domain(d3.extent(this._data, function(d) { return d[xName]; })).nice();
  svg.append('g').attr('class', 'axis')
                 .attr('transform', 'translate(0,' + (height/2) + ')')
                 .call(d3.axisBottom(x).tickPadding(10));
  var axisHeight = svg.select('.axis').node().getBoundingClientRect().height;
  // Set points
  svg.selectAll('.line').data(this._data)
     .enter().append('line').attr('id', function(d, i) { return 'pointLine'+i })
                            .classed('pointLine', true)
                            .attr('x1', function(d) { return x(d[xName]); })
                            .attr('x2', function(d) { return x(d[xName]); })
                            .attr('y1', height/2)
                            .attr('stroke-width', 1)
                            .attr('stroke', calculateColorOff)
                            .attr('opacity', 0.5)
                            .on('mouseover', generateMouseHandler(true))
                            .on('mouseout', generateMouseHandler(false));
  svg.selectAll('.circle').data(this._data)
     .enter().append('circle').attr('id', function(d, i) { return 'pointCircle'+i })
                              .attr('r', 5)
                              .attr('cx', function(d) { return x(d[xName]); })
                              .attr('cy', height/2)
                              .attr('fill', calculateColorOff)
                              .attr('opacity', 0.5)
                              .on('mouseover', generateMouseHandler(true))
                              .on('mouseout', generateMouseHandler(false));
  svg.selectAll('.text').data(this._data)
     .enter().append('text').text(returnText)
                            .attr('id', function(d, i) { return 'pointText'+i; })
                            .classed('pointText', true)
                            .attr('x', function(d) { return x(d[xName]); })
                            .style('font-size', 'small')
                            .style('font-weight', 'bold')
                            .style('text-anchor', 'middle')
                            .attr('fill', calculateColorOff)
                            .on('mouseover', generateMouseHandler(true))
                            .on('mouseout', generateMouseHandler(false));
  pointCircleOn = svg.append('circle').attr('id', function(d, i) { return 'pointCircleOn' })
                                      .attr('r', 8)
                                      .attr('cy', height/2)
                                      .style('display', 'none');
  // Calculate text positions.
  for (var i = 0; i < this._data.length; i++) {
    var datai = this._data[i];
    datai.x = x(datai[xName]);
    var recti = svg.select('#pointText'+i).node().getBoundingClientRect();
    datai.width = recti.width;
    datai.height = recti.height;
    var yupper = (height + axisHeight)/2 - axisHeight*0.75 - datai.height/2, ylower = (height + axisHeight)/2 + axisHeight*0.75 + datai.height/2;
    var k = 0;
    for (var resolved = false; !resolved;) {
      k++;
      var y = NaN;
      if (k%2) {
        if (yupper >= 0) {
          y = yupper;
        } else if (ylower + datai.height < height) {
          continue;
        }
      } else {
        if (ylower + datai.height < height) {
          y = ylower;
        } else if (yupper >= 0) {
          continue;
        }
      }
      if (isNaN(y)) {
        console.log("collision cannot be resolved");
        break;
      }
      resolved = true;
      for (var j = 0; j < i; j++) {
        var dataj = this._data[j];
        if ((Math.abs(dataj.x - datai.x) < (datai.width + dataj.width)*1.1/2) && (Math.abs(dataj.y - y) < (datai.height + dataj.height)/2)) {
          if (k%2) {
            yupper = y - (dataj.height + Math.max(y - dataj.y, 0));
          } else {
            ylower = y + (dataj.height + Math.max(dataj.y - y, 0));
          }
          resolved = false;
          break;
        }
      }
    }
    datai.y = k%2 ? yupper : ylower;
  }
  svg.selectAll('.pointText').attr('y', function(d) { return d.y; });
  svg.selectAll('.pointLine').attr('y2', function(d) { return d.y; });
};
function plot1d(container) {
  return new Plot1d(container);
}
