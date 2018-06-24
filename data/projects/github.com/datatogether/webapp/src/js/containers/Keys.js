import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import analytics from '../analytics';
import { selectLocalSessionUser, selectKeys } from '../selectors/session';
import { loadKeys, createKey, deleteKey } from '../actions/session';

import ValidInput from '../components/form/ValidInput';
import ValidTextarea from '../components/form/ValidTextarea';
import validateUser from '../validators/user';

import Spinner from '../components/Spinner';

class Keys extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      showErrors: false,
      saving: false,
      loading: true,
      newKey: {
        name: "",
        key: "",
      },
    };

    [
      'handleChange',
      'handleSave',
      'handleDelete',
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('keys');
    this.props.loadKeys();
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.keys && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  handleDelete(key) {
    this.props.deleteKey(key.name, key.sha256);
  }
  handleChange(name, value) {
    this.setState({ newKey: Object.assign({}, this.state.newKey, { [name]: value }) });
  }
  handleSave(e) {
    // const { validation, createKey } = this.props;
    e.preventDefault();

    // if (!validation.isValid) {
    //  if (!this.state.showErrors) {
    //    this.setState({ showErrors : true });
    //  }
    // } else {
    this.setState({ saving: true });
    this.props.createKey(this.state.newKey.name, this.state.newKey.key);
    // }
  }
  render() {
    const { newKey, showErrors, loading } = this.state;
    const { keys, validation } = this.props;

    if (loading) {
      return <Spinner />;
    }

    return (
      <div className="container">
        <div className="col-md-9 col-md-offset-3">
          <h3>Your SSH Keys:</h3>
          <ul>
            {keys.map((key, i) => {
              return (
                <li key={i}>
                  <p>{key.name}</p>
                  <small>SHA256:</small>
                  <p>{key.sha256}</p>
                  <button
                    className="btn btn-small btn-warning"
                    onClick={() => {
                      this.handleDelete.call(this, key);
                    }}
                  >Delete</button>
                </li>
              );
            })}
          </ul>
        </div>
        <div className="clear"></div>
        <div className="col-md-6 col-md-offset-3">
          <header>
            <h3>Add a Key:</h3>
            <ValidInput type="text" label="name" name="name" value={newKey.name} showError={showErrors} error={validation.name} onChange={this.handleChange} />
            <ValidTextarea label="paste key here" name="key" value={newKey.key} showError={showErrors} error={validation.description} onChange={this.handleChange} />
            <button className="btn btn-primary" disabled={showErrors} onClick={this.handleSave}>Add Key</button>
          </header>
        </div>
      </div>
    );
  }
}

Keys.propTypes = {
  // user: PropTypes.object,
  keys: PropTypes.array,
  validation: PropTypes.object,

  loadKeys: PropTypes.func.isRequired,
  createKey: PropTypes.func.isRequired,
  deleteKey: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  const user = selectLocalSessionUser(state);
  return {
    user,
    validation: validateUser(user),
    keys: selectKeys(state),
  };
}

export default connect(mapStateToProps, {
  loadKeys,
  createKey,
  deleteKey,
})(Keys);
