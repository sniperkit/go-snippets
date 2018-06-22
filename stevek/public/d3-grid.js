/*
 * Adapted from:
 * http://bl.ocks.org/herrstucki/5684816
 */
(function() {
    d3.layout.grid = function() {
        var mode = "equal",
            layout = _distributeEqually,
            x = d3.scale.ordinal(),
            y = d3.scale.ordinal(),
            size = [1, 1],
            actualSize = [0, 0],
            nodeSize = false,
            cols, rows;

        function grid(nodes) {
            return layout(nodes);
        }

        function _distributeEqually(nodes) {
            var i = -1,
                n = nodes.length,
                _cols = cols ? cols : 0,
                _rows = rows ? rows : 0,
                col, row;

            if (_rows && !_cols) {
                _cols = Math.ceil(n / _rows)
            } else {
                _cols || (_cols = Math.ceil(Math.sqrt(n)));
                _rows || (_rows = Math.ceil(n / _cols));
            }

            x.domain(d3.range(_cols)).rangePoints([0, size[0]]);
            y.domain(d3.range(_rows)).rangePoints([0, size[1]]);
            actualSize[0] = x(1);
            actualSize[1] = y(1);

            while (++i < n) {
                col = i % _cols;
                row = Math.floor(i / _cols);
                nodes[i].x = x(col);
                nodes[i].y = y(row);
            }

            return nodes;
        }

        grid.size = function(value) {
            if (!arguments.length) return nodeSize ? actualSize : size;
            actualSize = [0, 0];
            nodeSize = (size = value) == null;
            return grid;
        }

        grid.nodeSize = function(value) {
            if (!arguments.length) return nodeSize ? size : actualSize;
            actualSize = [0, 0];
            nodeSize = (size = value) != null;
            return grid;
        }

        grid.rows = function(value) {
            if (!arguments.length) return rows;
            rows = value;
            return grid;
        }

        grid.cols = function(value) {
            if (!arguments.length) return cols;
            cols = value;
            return grid;
        }

        grid.points = function() {
            return grid;
        }

        return grid;
    };
})();
