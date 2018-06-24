import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import analytics from '../analytics';
import { selectSessionUser } from '../selectors/session';
import { selectUncrawlables } from '../selectors/uncrawlables';
import { loadUncrawlables } from '../actions/uncrawlables';

import List from '../components/List';
import UncrawlableItem from '../components/item/UncrawlableItem';
import Spinner from '../components/Spinner';

class Uncrawlables extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
    };

    [].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('uncrawlables');
    // this.props.loadUncrawlables();
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.uncrawlables && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  render() {
    const { loading } = this.state;
    const { uncrawlables, user } = this.props;

    if (loading) {
      return <Spinner />;
    }

    return (
      <div id="uncrawlables" className="page">
        <div className="container">
          <header className="row">
            <div className="col-md-12">
              <hr className="blue" />
              { user && (<Link className="btn bg-blue white right" to="/uncrawlables/new">New Uncrawlable</Link>)}
              <h1 className="blue">Uncrawlables:</h1>
              <p>Uncrawlables are urls that can't be archived without custom code</p>
              <hr className="blue" />
            </div>
          </header>
          <div className="row">
            <br />
            <List data={uncrawlables} component={UncrawlableItem} />
          </div>
        </div>
      </div>
    );
  }
}

Uncrawlables.propTypes = {
  // user: PropTypes.object,
  uncrawlables: PropTypes.array,
  loadUncrawlables: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    user: selectSessionUser(state),
    uncrawlables: selectUncrawlables(state),
  };
}

export default connect(mapStateToProps, {
  loadUncrawlables,
})(Uncrawlables);
