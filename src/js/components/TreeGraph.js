import React, { PropTypes, Children, cloneElement } from 'react';
import { hierarchy, tree } from 'd3-hierarchy';

import PureRenderComponent from './PureRenderComponent';
// import * as forceUtils from '../utils/d3-force';
// import * as rafUtils from '../utils/raf';
// import ZoomableSVGGroup from './ZoomableSVGGroup';
// import simulationPropTypes, { DEFAULT_SIMULATION_PROPS } from '../propTypes/simulation';

import { maxDepthCount, widestDepthCount } from '../utils/tree';

// export function isNode(child) {
//   return child.props && child.props.node;
// }

// export function isLink(child) {
//   return child.props && child.props.link;
// }

export default class TreeGraph extends PureRenderComponent {
  render() {
    const { data, layout, onSelectNode, colorNode, selectedNodes } = this.props;
    const { width, height } = layout;
    const hasSelection = (selectedNodes && selectedNodes.length > 0);

    if (!data) {
      return null;
    }

    const root = hierarchy(data);

    function nodeIsSelected(node) {
      return selectedNodes.findIndex((n) => node.data.id == n.data.id) >= 0;
    }

    let w = maxDepthCount(data) * 300;
    if (w < width) {
      w = width;
    }
    // determine how tall window needs to be
    let h = widestDepthCount(data) * 20;
    if (h < height) {
      h = height;
    }
    // const t = tree().size([height - Math.floor(height * 0.2), width - Math.floor(width * 0.2)])(root);
    const t = tree().size([h, w - Math.floor(w * 0.2)])(root);

    return (
      <svg style={{
        position: "absolute",
        width: w + 300,
        height: h + 300,
        top: layout.top,
        left: layout.left,
      }}>
        <g transform={`translate(${Math.floor(w * 0.1)}, ${Math.floor(height * 0.1)})`}>
          <g className="links">
            {t.descendants().slice(1).map((d) => {
              return (
                <path
                  key={d.data.id}
                  className="link"
                  fill="none"
                  stroke="#DDDDDD"
                  strokeWidth="1"
                  d={
                      "M" + d.y + "," + d.x
                      + "C" + (d.y + d.parent.y) / 2 + "," + d.x
                      + " " + (d.y + d.parent.y) / 2 + "," + d.parent.x
                      + " " + d.parent.y + "," + d.parent.x
                    }
                />);
            })}
          </g>
          <g className="nodes">
            {t.descendants().map((d) => {
              const leaf = (!d.data.numChildren || d.data.numChildren == 0)
              return (
                <g 
                  key={d.data.id}
                  className={ d.children ? "node node--internal" : "node node--leaf"}
                  transform={`translate(${d.y},${d.x})`}
                >
                <text 
                  dy="3"
                  x={d.children ? -8 : 8}
                  textAnchor={d.children ? "end" : "start"}
                  onClick={(e) => { onSelectNode.call(this, e, d) }}
                  style={{ opacity : hasSelection ? nodeIsSelected(d) ? 1.0 : 0.25 : 1.0}}
                >{d.data.name}</text>
                <circle 
                  r="6"
                  fill={colorNode(d.data)}
                  strokeWidth={ leaf ? 3 : 0}
                  stroke={ leaf ? "#999999" : "none" }
                  onClick={(e) => { onSelectNode.call(this, e, d) }}
                  style={{ opacity : hasSelection ? nodeIsSelected(d) ? 1.0 : 0.25 : 1.0}}
                ></circle>
              </g>)
            })}
          </g>
        </g>
      </svg>
    );
  }
}

TreeGraph.propTypes = {
  layout: PropTypes.object,
  data: PropTypes.object,
  rootNode: PropTypes.object,
  selectedNodes: PropTypes.array,
  onSelectNode: PropTypes.func,
  colorNode: PropTypes.func.isRequired,
};

TreeGraph.defaultProps = {
  layout: {
    width: 300,
    height: 300,
    t: 0,
    l: 0,
  },
};
