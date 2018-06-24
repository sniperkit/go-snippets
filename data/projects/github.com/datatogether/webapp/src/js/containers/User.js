import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';

import analytics from '../analytics';
import { loadUserByUsername, loadUserCollections } from '../actions/user';
import { logoutUser, loadKeys } from '../actions/session';
import { loadMetadataByKey } from '../actions/metadata';
import { loadCommunityUsers } from '../actions/communities';

import { selectSessionUser } from '../selectors/session';
import { selectUserByUsername, selectCommunityUsers } from '../selectors/user';
import { selectCollectionsByKey } from '../selectors/collections';
// import { selectDefaultKeyId } from '../selectors/keys';

import NotFound from '../components/NotFound';
import List from '../components/List';
import CollectionItem from '../components/item/CollectionItem';
import UserItem from '../components/item/UserItem';

class User extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: !this.props.user,
    };

    [
      "handleLogout",
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('user');
    if (this.props.user) {
      if (this.props.user.type == "community") {
        this.props.loadCommunityUsers(this.props.user.id, 1, 25);
      }
      if (this.props.user.currentKey) {
        // this.props.loadMetadataByKey(this.props.user.currentKey);
        this.props.loadUserCollections(this.props.user.currentKey, 1, 25);
      }
    } else {
      this.props.loadUserByUsername(this.props.username);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.username != this.props.username) {
      nextProps.loadUserByUsername(nextProps.username);
      this.setState({ loading: true });
    } else if (nextProps.user && this.state.loading) {
      if (nextProps.user) {
        if (nextProps.user.type == "community") {
          this.props.loadCommunityUsers(nextProps.user.id, 1, 25);
        }
        if (nextProps.user.currentKey) {
          // this.props.loadMetadataByKey(nextProps.user.currentKey);
          this.props.loadUserCollections(nextProps.user.currentKey, 1, 25);
        }
      }
      this.setState({ loading: false });
    }
  }

  handleLogout() {
    this.props.logoutUser();
  }

  render() {
    const { user, communityUsers, permissions, collections } = this.props;
    if (!user) {
      return <NotFound />;
    }

    return (
      <div id="user" className="user page">
        <header className="colorized">
          <div className="container">
            <div className="row">
              <div className="col-md-12">
                <hr className="" />
              </div>
              <div className="col-md-2">
                {user.profileUrl &&
                  <div className="profile_photo">
                    <img src={user.profileUrl} />
                  </div>
                }
              </div>
              <div className="col-md-10">
                <label className="label">{user.type || "User"}</label>
                <Link className="" to={`/${user.username}`}>
                  <h1 className="title">{user.username}</h1>
                </Link>
                { permissions.edit ? <Link to="/settings" >settings </Link> : undefined }
                { permissions.edit ? <a onClick={this.handleLogout}>logout</a> : undefined }
                <p>{ user.description }</p>
              </div>
            </div>
          </div>
        </header>
        <div className="container">
          <div className="row">
            <div className="col-md-9">
              <h4 className="title">Collections:</h4>
              <hr />
              <List data={collections} component={CollectionItem} />
            </div>
            {user.type == "community" && communityUsers &&
            <div className="col-md-3">
              <h4 className="title">People:</h4>
              <hr />
              <List data={communityUsers} component={UserItem} />
            </div>}
          </div>
        </div>
      </div>
    );
  }
}

User.propTypes = {
  username: PropTypes.string.isRequired,
  user: PropTypes.object,
  permissions: PropTypes.object.isRequired,
  loadUserByUsername: PropTypes.func.isRequired,
  logoutUser: PropTypes.func.isRequired,
};

User.defaultProps = {
};

function mapStateToProps(state, ownProps) {
  const username = ownProps.params.user;
  const user = selectUserByUsername(state, username);

  let permissions = {
    edit: false,
    del: false,
  };

  const session = selectSessionUser(state);
  if (session && session.username == username) {
    permissions.edit = true;
    permissions.del = true;
  }

  return Object.assign({
    username,
    user,
    communityUsers: user && user.type == "community" ? selectCommunityUsers(state, user.id) : [],
    permissions,
    collections: user ? selectCollectionsByKey(state, user.currentKey) : [],
  }, ownProps);
}

export default connect(mapStateToProps, {
  loadUserByUsername,
  loadCommunityUsers,
  logoutUser,
  loadKeys,
  loadMetadataByKey,
  loadUserCollections,
})(User);
