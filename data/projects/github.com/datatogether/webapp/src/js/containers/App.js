/* global window */
import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { debounce } from 'lodash';

import analytics from '../analytics';
import * as socket from '../middleware/socket';

import {
  toggleMenu,
  hideMenu,
  setErrorMessage,
  resetMessage,
  resetErrorMessage,
  showModal,
  hideModal,
} from '../actions/app';
import { layoutResize } from '../actions/layout';
import { loadSessionUser } from '../actions/session';
import { selectSessionUser } from '../selectors/session';

import Navbar from '../components/Navbar';
import MainMenu from '../components/MainMenu';

const BETA_SIGNUP_MODAL = 'BETA_SIGNUP_MODAL';

class App extends Component {
  constructor(props) {
    super(props);

    [
      "handleChange",
      "handleDismissErrorMessage",
      "handleDismissStatusMessage",
      "handleHideMenu",
      "handleMenuToggle",
      "handleGimmieInvite",
      "handleHideModal",
      "modal",
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    this.props.dispatch(loadSessionUser());

    // attempt to connect websocket, passing in the dispatcher
    // TODO - move this up the app boostrapping stack
    socket.connect(this.props.dispatch, 6500).then(() => {
      // socket connected.
    }).catch((err) => {
      console.log("conn error:", err);
      this.props.dispatch(setErrorMessage("connection error"));
    });

    this._oldOnbeforeunload = window.onbeforeunload;
    window.onbeforeunload = function() {
      socket.disconnect();
    };

    this._oldResize = window.onresize;
    // debounce device resizing to not be a jerk on resize
    window.onresize = debounce(() => {
      this.props.dispatch(layoutResize(window.innerWidth, window.innerHeight));
    }, 250);

    // initial call to make things not crazy
    this.props.dispatch(layoutResize(window.innerWidth, window.innerHeight));
  }

  componentWillReceiveProps(nextProps) {
    if (!this.props.user && nextProps.user) {
      analytics.identify(nextProps.user.id, {
        username: nextProps.user.username,
      });
    }
  }

  componentWillUnmount() {
    // restore old onResize if one exists b/c we're good citizens like that
    window.onresize = this._oldResize;
    window.onbeforeunload = this._oldOnbeforeunload;

    // teardown websocket connection
    socket.disconnect();
  }

  handleChange(nextValue) {
    browserHistory.push(`/${nextValue}`);
  }
  handleDismissErrorMessage(e) {
    this.props.dispatch(resetErrorMessage());
    e.preventDefault();
  }

  handleDismissStatusMessage(e) {
    this.props.dispatch(resetMessage());
    e.preventDefault();
  }

  handleMenuToggle(e) {
    e.stopPropagation();
    this.props.dispatch(toggleMenu());
  }
  handleHideMenu() {
    if (this.props.showMenu) {
      this.props.dispatch(hideMenu());
    }
  }

  handleGimmieInvite() {
    this.props.dispatch(hideMenu());
    this.props.dispatch(showModal(BETA_SIGNUP_MODAL, this));
  }
  handleHideModal() {
    this.props.dispatch(hideModal());
  }


  /* app implements the modal pattern as well as using it */
  modal(name) {
    switch (name) {
      // case BETA_SIGNUP_MODAL:
      //   return <BetaSignup onSaved={this.handleHideModal} onCancelled={this.handleHideModal} />
      default:
        return undefined;
    }
  }

  /**
   * presenting modals is easy! fun even! yay! in your component of choice
   * import showModal from actions/app, and call it with the signature
   * ("name", [component that will present modal], [data object])
   * it's expected that the element that presents the modal will have a method "modal",
   * that will return either a react element or undefined
   * Whatever it gives back will be presented
   */
  renderModal() {
    if (this.props.modal) {
      return (
        <div id="modal-wrap">
          <div className="modal dialog" tabIndex="-1" role="dialog">
            {this.props.modal.element ? this.props.modal.element.modal(this.props.modal.name, this.props.modal.data) : undefined}
          </div>
        </div>
      );
    }

    return undefined;
  }

  renderErrorMessage() {
    const { errorMessage } = this.props;
    if (!errorMessage) {
      return null;
    }

    return (
      <div className="alert container red">
        <div className="row">
          <div className="col-md-12">
            <p className="message">{errorMessage}</p>
            <a className="dismiss" onClick={this.handleDismissErrorMessage}>Dismiss</a>
          </div>
        </div>
      </div>
    );
  }

  renderStatusMessage() {
    const { statusMessage } = this.props;
    if (!statusMessage) {
      return null;
    }

    return (
      <div className="alert alert-success" role="alert">
        <b>{statusMessage}</b>
        <a onClick={this.handleDismissStatusMessage}>Dismiss</a>
      </div>
    );
  }

  render() {
    const { children, user, showMenu, layout } = this.props;
    return (
      <div id="app" onClick={this.handleHideMenu}>
        <Navbar
          user={user}
          style={Object.assign({
            position: "absolute",
          }, layout.navbar)}
          onToggleMenu={this.handleMenuToggle}
          onGimmieInvite={this.handleGimmieInvite}
        />
        <MainMenu user={user} show={showMenu} />
        {this.renderErrorMessage()}
        {this.renderStatusMessage()}
        {children}
        {this.renderModal()}
      </div>
    );
  }
}

App.propTypes = {
  // Injected by React Redux
  errorMessage: PropTypes.string,
  // message: PropTypes.string,
  // inputValue: PropTypes.string.isRequired,
  // Injected by React Router
  children: PropTypes.node,
  user: PropTypes.object,

  // resetErrorMessage: PropTypes.func.isRequired,
  // resetMessage: PropTypes.func.isRequired,
  // loadSessionUser: PropTypes.func.isRequired,
  // hideMenu: PropTypes.func.isRequired
};

function mapStateToProps(state, ownProps) {
  return {
    statusMessage: state.statusMessage,
    errorMessage: state.errorMessage,
    inputValue: ownProps.location.pathname.substring(1),
    user: selectSessionUser(state),
    showMenu: state.app.showMenu,
    layout: state.layout,

    modal: state.app.modal,
  };
}

// function createHandlers(dispatch) {
//   let onClick = function(node, data) {
//     dispatch(actions.nodeClicked(data))
//   };

//   return {
//     resetMessage : () => dispatch(resetMessage),
//     resetErrorMessage,
//     loadSessionUser,
//     layoutResize,
//     hideMenu,
//     showModal,
//     hideModal
//   }
// }


export default connect(mapStateToProps)(App);
