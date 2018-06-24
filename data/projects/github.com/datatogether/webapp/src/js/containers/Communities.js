import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import analytics from '../analytics';
import { selectCommunities } from '../selectors/user';
import { loadUsers } from '../actions/user';

import List from '../components/List';
import UserItem from '../components/item/UserItem';
import Spinner from '../components/Spinner';

class Communities extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
    };

    [].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('communities');
    this.props.loadUsers("community", 1, 25);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.communities && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  render() {
    const { loading } = this.state;
    const { communities } = this.props;

    if (loading) {
      return <Spinner />;
    }

    return (
      <div id="communities" className="community page">
        <header className="colorized">
          <div className="container">
            <div className="row">
              <div className="col-md-12">
                <hr />
                <h1>Communities:</h1>
              </div>
              <div className="col-md-12">
                <hr className="" />
              </div>
            </div>
          </div>
        </header>
        <div className="container">
          <div className="row">
            <br />
            <List data={communities} component={UserItem} />
          </div>
        </div>
      </div>
    );
  }
}

Communities.propTypes = {
  // user: PropTypes.object,
  communities: PropTypes.array,
  loadUsers: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    communities: selectCommunities(state),
  };
}

export default connect(mapStateToProps, {
  loadUsers,
})(Communities);
