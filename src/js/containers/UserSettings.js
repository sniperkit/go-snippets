import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import analytics from '../analytics';
import { selectAvailableUsers, selectSessionUser, selectLocalSessionUser } from '../selectors/session';
import { selectUserByUsername, selectLocalUser } from '../selectors/user';
import { editUser, updateUser, saveUser } from '../actions/user';

import ValidInput from '../components/form/ValidInput';
import ValidTextarea from '../components/form/ValidTextarea';
import validateUser from '../validators/user';

import Spinner from '../components/Spinner';
import NotFound from '../components/NotFound';

class UserSettings extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      showErrors: false,
      saving: false,
    };

    [
      'handleChange',
      'handleSave',
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('user_settings');
    if (this.props.savedUser) {
      this.props.editUser(savedUser.id);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.user && this.state.loading) {
      this.setState({ loading: false });
    }
    if (!this.props.savedUser && nextProps.savedUser) {
      this.props.editUser(nextProps.savedUser.id);
    }
  }

  handleChange(name, value) {
    const user = Object.assign({}, this.props.user);
    user[name] = value;
    this.props.updateUser(user);
  }
  handleSave(e) {
    const { user, validation } = this.props;
    e.preventDefault();

    if (!validation.isValid) {
      if (!this.state.showErrors) {
        this.setState({ showErrors: true });
      }
    } else {
      this.setState({ saving: true });
      this.props.saveUser(user);
    }
  }
  render() {
    const { showErrors, loading } = this.state;
    const { savedUser, user, validation, editableUsers, editable } = this.props;

    if (loading) {
      return <Spinner />;
    }

    if (!user || !editable) {
      return <NotFound />
    }

    return (
      <div id="userSettings" className="page">
        <div className="container">
          {editableUsers && <div className="col-md-3">
            <ul>
              {editableUsers.map((u,i) => {
                return (
                  <li key={i}><Link to={`/users/${u.username}/settings`}>{u.username}</Link></li>
                );
              })}
            </ul>
          </div>}
          <div className="col-md-6 col-md-offset-3">
            <header>
              <h3>{savedUser.username} Settings</h3>
              <ValidInput type="text" label="name" name="name" value={user.name} showError={showErrors} error={validation.name} onChange={this.handleChange} />
              <ValidInput type="text" label="username" name="username" value={user.username} showError={showErrors} error={validation.username} onChange={this.handleChange} />
              <ValidInput type="text" label="profile photo url" name="profileUrl" value={user.profileUrl} showError={showErrors} error={validation.profileUrl} onChange={this.handleChange} />
              <ValidInput type="text" label="profile thumb url" name="thumbUrl" value={user.thumbUrl} showError={showErrors} error={validation.thumbUrl} onChange={this.handleChange} />
              <ValidInput type="text" label="email" name="email" value={user.email} showError={showErrors} error={validation.email} onChange={this.handleChange} />
              <ValidTextarea label="bio/description/about" name="description" value={user.description} showError={showErrors} error={validation.description} onChange={this.handleChange} />
              <button className="btn btn-primary" disabled={(!validation.isValid && showErrors)} onClick={this.handleSave}>Update</button>
            </header>
          </div>
        </div>
      </div>
    );
  }
}

UserSettings.propTypes = {
  user: PropTypes.object,
  validation: PropTypes.object,
  editableUsers: PropTypes.array,

  editUser: PropTypes.func.isRequired,
  updateUser: PropTypes.func.isRequired,
  saveUser: PropTypes.func.isRequired,
};

UserSettings.defaultProps = {

};

function mapStateToProps(state, ownProps) {
  const editableUsers = selectAvailableUsers(state);
  const username = ownProps.params.user;
  const savedUser = selectUserByUsername(state, username);

  if (!selectSessionUser(state)) {
    return ({
      editable: false,
      editableUsers,
    });
  }

  if (savedUser) {
    console.log(editableUsers);
    console.log(editableUsers.findIndex((u) => savedUser.id == u.id));
  }

  const editable = savedUser ? editableUsers.findIndex((u) => savedUser.id == u.id) >= 0 : false;
  const user = selectLocalUser(state, savedUser ? savedUser.id : "");

  return {
    editable,
    editableUsers,
    savedUser,
    user,
    validation: validateUser(user),
  };
}

export default connect(mapStateToProps, {
  editUser,
  updateUser,
  saveUser,
})(UserSettings);
