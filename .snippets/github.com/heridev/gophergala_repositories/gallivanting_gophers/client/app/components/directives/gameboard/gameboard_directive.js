angular.module('gg.components.gameboard')
.directive('ggGameboard', function() {
  return {
    restrict: 'E',
    scope: {
      updateData: '='
    },
    controller: 'gameboardCtrl',
    templateUrl: 'components/directives/gameboard/gameboard.html',
    replace: true,
    link: function (scope, element, attrs) {
      var board = d3.select(element[0]);
      var svg = board.append("svg")
        .attr("width", 840)
        .attr("height", 840)
        .attr("class", "gameboard");
      var currentSelected;

      scope.$watch('updateData', function(newValue, oldValue) {
        if (newValue != oldValue && typeof newValue != 'undefined') {
          scope.updateBoard(scope.updateData);
        }
      });

      scope.selectGopher = function(g) {
        var elem = angular.element(g);
        elem.attr('r', 25);
        currentSelected = elem.attr('id').split('-')[1];
        scope.gopherSelected(currentSelected);
      };

      scope.updateBoard = function(data) {
        for (var y=0;y<data.tiles.length;y++) {
          var vertical = data.tiles[y] >>> 16;
          var horizontal = data.tiles[y] << 16;
          horizontal = horizontal >>> 16;
          for (var x=0;x<16;x++) {
            var vcalc = vertical & (1 << (15-x));
            var hcalc = horizontal & (1 << (15-x));
            if (vcalc > 0) {
              angular.element(document.getElementById('border-vertical-' + x + '-' + y)).css('fill', '#000');
            } else {
              angular.element(document.getElementById('border-vertical-' + x + '-' + y)).css('fill', 'rgba(240, 240, 240, .75)');
            }
            if (hcalc > 0) {
              angular.element(document.getElementById('border-horizontal-' + x + '-' + y)).css('fill', '#000');
            } else {
              angular.element(document.getElementById('border-horizontal-' + x + '-' + y)).css('fill', 'rgba(240, 240, 240, .75)');
            }
          }
        }

        for (var g=0;g<data.goals.length;g++) {
          var goal = angular.element(document.getElementById('goal-' + g));
          if (typeof goal != 'undefined') {
            goal.remove();
          }
        }

        for (var i=0;i<data.goals.length;i++) {
          var block = data.goals[i].Location;
          var row  = Math.floor(block / 16);
          var column = block % 16;

          svg.append("svg:image")
          .attr("id", "goal-" + i)
          .attr("x", ((column * 52)))
          .attr("y", ((row * 52)))
          .attr("width", 50)
          .attr("height", 50)
          .attr("xlink:href", "imgs/goal-" + i +".svg")
        }

        for (var g=0;g<data.gophers.length;g++) {
          var goph = angular.element(document.getElementById('gopher-' + g));
          if (typeof goph != 'undefined') {
            goph.remove();
          }
        }

        for (var c=0;c<data.gophers.length;c++) {
          var gblock = data.gophers[c].Location;
          var grow = Math.floor(gblock / 16);
          var gcolumn = gblock % 16;
          var gcolor;

          switch (data.gophers[c].Type) {
            case 1: gcolor = 'red';
            break;
            case 2: gcolor = 'green';
            break;
            case 3: gcolor = 'purple';
            break;
            case 4: gcolor = 'yellow';
            break;
            default: gcolor = 'blue';
          }

          svg.append("svg:image")
            .attr("id", "gopher-" + c)
            .attr("x", ((gcolumn * 52)))
            .attr("y", ((grow * 52)))
            .attr("width", 50)
            .attr("height", 50)
            .attr("xlink:href", "imgs/gopher-" + gcolor +".svg")
            .on("click", function() { scope.selectGopher(this); });
        }
      };

      var _createBoard = function() {
        for (var y=0;y<16;y++) {
          for (var x=0;x<16;x++) {
            svg.append("rect")
            .attr("id", "block-" + x + "-" + y)
            .attr("x", x * 52)
            .attr("y", y * 52)
            .attr("width", 50)
            .attr("height", 50)
            .attr("class", "gamesquare");

            if (x < 15) {
              svg.append("rect")
              .attr("id", "border-vertical-" + x + "-" + y)
              .attr("x", (x * 52) + 50)
              .attr("y", y * 52)
              .attr("width", 2)
              .attr("height", 50);
            }
            if (y < 15) {
              svg.append("rect")
              .attr("id", "border-horizontal-" + x + "-" + y)
              .attr("x", x * 52)
              .attr("y", (y * 52) + 50)
              .attr("width", 50)
              .attr("height", 2)
              .attr("class", "squareboarder");
            }
          }
        }
      };

      _createBoard();
    }


  };
});
