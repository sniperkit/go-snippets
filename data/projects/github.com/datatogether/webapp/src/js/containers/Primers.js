import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import analytics from '../analytics';
import { selectPrimers } from '../selectors/primers';
import { loadPrimers } from '../actions/primers';

import List from '../components/List';
import PrimerItem from '../components/item/PrimerItem';
import Spinner from '../components/Spinner';

class Primers extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
    };

    [].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('primers');
    this.props.loadPrimers(1, 25, true);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.primers && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  render() {
    const { loading } = this.state;
    const { primers } = this.props;

    if (loading) {
      return <Spinner />;
    }

    return (
      <div id="primers" className="primer page">
         <header className="container">
          <div className="row">
            <div className="col-md-12">
              <hr />
              <h1>Primers:</h1>
            </div>
            <div className="col-md-6">
              <p>In order to make our archiving efforts as thorough and systematic as possible,
              we use Agency Archiving Primers to identify key programs, 
              datasets, and documents that are vulnerable to change and loss. Primers are composed of <i>sources</i>, 
              which each specify a url as a starting point for archiving.</p>
            </div>
            <div className="col-md-12">
              <hr className="" />
            </div>
          </div>
        </header>
        <div className="container">
          <div className="row">
            <br />
            <List data={primers} component={PrimerItem} />
          </div>
        </div>
      </div>
    );
  }
}

Primers.propTypes = {
  // user: PropTypes.object,
  primers: PropTypes.array,
  loadPrimers: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    primers: selectPrimers(state),
  };
}

export default connect(mapStateToProps, {
  loadPrimers,
})(Primers);
