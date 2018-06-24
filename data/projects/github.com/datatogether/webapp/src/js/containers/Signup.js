import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';

import analytics from '../analytics';
import { signup, SESSION_SIGNUP_FAILURE } from '../actions/session';
import { selectSessionUser } from '../selectors/session';
import ValidInput from '../components/form/ValidInput';
import validateUser from '../validators/user';

// TODO - add validation logic
class Signup extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: "",
      name: "",
      email: "",
      password: "",
      showErrors: false,
    };

    ['handleChange', 'handleSignupSubmit'].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('signup');
    if (this.props.user != null) {
      browserHistory.push('/');
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.user != null) {
      browserHistory.push('/');
    }
  }

  handleSignupSubmit(e) {
    e.preventDefault();
    const validation = validateUser(this.state);
    if (validation.isValid) {
      this.setState({ loading: true });
      this.props.signup(this.state).then((action) => {
        if (action.type == SESSION_SIGNUP_FAILURE) {
          this.setState({ loading: false });
        }
      });
    } else if (!this.state.showErrors) {
      this.setState({ showErrors: true });
    }
  }

  handleChange(name, value) {
    const inv = Object.assign({}, this.state, { [name]: value });
    this.setState(inv);
  }

  render() {
    const { username, email, password, showErrors, loading } = this.state;
    const validation = validateUser(this.state);

    return (
      <div id="signup" className="page">
        <div className="container">
          <form onSubmit={this.handleSignupSubmit} className="form col-md-4 offset-md-4">
            <h3>Signup:</h3>
            <ValidInput type="text" label="Username" name="username" value={username} error={validation.username} showError={showErrors} onChange={this.handleChange} />
            <ValidInput type="text" label="Email" name="email" value={email} error={validation.email} showError={showErrors} onChange={this.handleChange} />
            <ValidInput type="password" label="Password" name="password" value={password} error={validation.password} showError={showErrors} onChange={this.handleChange} />
            <input className="signup btn btn-standard" disabled={loading} type="submit" value="signup" />
          </form>
        </div>
      </div>
    );
  }
}

Signup.propTypes = {
  signup: PropTypes.func.isRequired,
};

Signup.defaultProps = {

};

function mapStateToProps(state) {
  return {
    user: selectSessionUser(state),
  };
}

export default connect(mapStateToProps, {
  signup,
})(Signup);
