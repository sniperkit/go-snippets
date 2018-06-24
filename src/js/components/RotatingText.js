import React, { PropTypes } from 'react';


export default class RotatingText extends React.Component {
  constructor(props) {
    super(props);
    
    this.state = {
      index: 0,
      fadeState: "display",
    };

    [
      "handleStartFade",
      "handleStopFade",
    ].forEach((m) => this[m] = this[m].bind(this));
  }

  componentWillMount() {
    this.timer = setInterval(this.handleStartFade, this.props.displayTime + this.props.transitionTime);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.displayTime != this.props.displayTime) {
      this.timer = setInterval(this.handleStartFade, this.props.displayTime + this.props.transitionTime);
    }
  }

  handleStartFade() {
    const { index } = this.state;
    this.setState({ fadeState: "fade-out" });

    setTimeout(() => {
      this.setState({ 
        index: (index + 1 == this.props.text.length) ? 0 : index + 1,
        fadeState: "fade-in",
      });
      setTimeout(this.handleStopFade, this.props.transitionTime / 2);
    }, this.props.transitionTime / 2);
  }

  handleStopFade() {
    this.setState({ fadeState: "display" });
  }

  style(props, fadeState) {
    return {
      color: props.color,
      minWidth: props.minWidth,
      transition : `all ${this.props.transitionTime / 1000}s`,
      opacity: (fadeState == "fade-out") ? 0 : 1,
    }
  }

  render() {
    const { text } = this.props;
    const { index, fadeState } = this.state;

    return (
      <span className={`rotating text ${fadeState}`} style={this.style(this.props, fadeState)}>
        <span>{text[index]}</span>
      </span>
    );
  }
}

RotatingText.propTypes = {
  text: PropTypes.array.isRequired,
  minWidth: PropTypes.number,
  displayTime: PropTypes.number,
  transitionTime: PropTypes.number,
  color: PropTypes.string,
};

RotatingText.defaultProps = {
  minWidth: 200,
  displayTime: 1000,
  transitionTime: 500,
  color: "#FFFFFF",
}
