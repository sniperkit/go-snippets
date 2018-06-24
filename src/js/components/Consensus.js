import React, { PropTypes } from 'react';

export default class Consensus extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};

    [
      "handleKeyClick",
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }


  handleKeyClick(key) {
    this.setState({
      [key]: !this.state[key],
    });
  }

  renderAllKeys(key) {
    if (!key) { return undefined; }
    return (
      <div>
        {key.map((val, i) => {
          return <p key={i}>{val || <i>blank</i>}</p>;
        })}
      </div>
    );
  }

  render() {
    const { data } = this.props;

    return (
      <div className="consensus">
        {Object.keys(data).map((key) => {
          return (
            <div key={key}>
              <label className="key" onClick={this.handleKeyClick.bind(this, key)}><span className="yellow meta label">{key}</span> <i style={{ opacity: 0.4 }}>{data[key].length > 1 ? data[key].length : "" }</i></label>
              { this.state[key] ? this.renderAllKeys(data[key]) : <p>{data[key][0] ? data[key][0] : <i>blank</i>}</p> }
            </div>
          );
        })}
      </div>
    );
  }
}

Consensus.propTypes = {
  data: PropTypes.object.isRequired,
};

Consensus.defaultProps = {
  onClickKey: () => {},
};

