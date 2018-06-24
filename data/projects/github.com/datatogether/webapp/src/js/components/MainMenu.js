import React, { PropTypes } from 'react'
import { Link } from 'react-router'

export default class MainMenu extends React.Component {
  onClick (e) {
    // need to prevent any clicks on mainmenu from bubbling up
    // & hiding the menu
    e.stopPropagation()
  }

  userLinks (user) {
    if (!user) {
      return <Link to='/login'>Login</Link>
    }

    return (
      <div>
        <Link to={`/user/${user.username}`}>{user.username}</Link>
        <a href='/docs'>Docs</a>
      </div>
    )
  }

  render () {
    const { user, show } = this.props
    return (
      <div id='main_menu' onClick={this.onClick} className={show ? 'show' : 'hide'}>
        <Link className='blue' to='/'>Home</Link>
        <Link className='blue' to='/add-dataset'>Add a Dataset</Link>
        <Link className='blue' to='/communities'>Communities</Link>
        <Link className='blue' to='/collections'>Collections</Link>
        {
          !user &&
            (<div>
              <Link className='yellow' to='/login'>Login</Link>
              <Link className='yellow' to='/signup'>Signup</Link>
            </div>)
        }
      </div>
    )
  }
}

MainMenu.propTypes = {
  user: PropTypes.oneOfType([
    PropTypes.object,
    PropTypes.null]),
  show: PropTypes.bool
}

MainMenu.defaultProps = {
}
