/* globals __BUILD__ */
import React, { PropTypes } from 'react';
import { Link } from 'react-router';

const Navbar = ({ user, style, onToggleMenu }) => {
  return (
    <div id="navbar" style={style}>
      <div className="container">
        <div className="row">
          <div className="col-md-6 col-sm-6 col-xs-6">
            <Link id="logotype" href={__BUILD__.BASE_URL}>[ data, together ]</Link>
          </div>
          <div className="menu col-md-4 offset-md-2 col-sm-6">
            {user && <Link to={`/users/${user.username}`}>{user.username}</Link>}
            <a className="green" onClick={onToggleMenu}>Menu</a>
          </div>
        </div>
      </div>
    </div>
  );
};

Navbar.propTypes = {
  user: PropTypes.oneOfType([
    React.PropTypes.object,
    React.PropTypes.null]),
  onToggleMenu: PropTypes.func.isRequired,
};

Navbar.defaultProps = {
};

export default Navbar;
