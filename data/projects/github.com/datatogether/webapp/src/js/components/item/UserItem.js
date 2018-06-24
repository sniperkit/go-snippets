import React from 'react';
import { Link } from 'react-router';

const UserItem = ({ data, onSelect }) => {
  const imgUrl = data.thumbUrl || data.profileUrl;

  return (
    <div className="user item">
      { imgUrl &&
        <Link to={`/users/${data.username}`} className="profile-photo left">
          <img src={imgUrl} />
        </Link>
      }
      <h4 onClick={onSelect}>
        <Link className="namespace" to={`/users/${data.username}`}>{data.username}</Link>
      </h4>
      <div className="clear"></div>
    </div>
  );
};

UserItem.propTypes = {
  data: React.PropTypes.object.isRequired,
  onSelect: React.PropTypes.func.isRequired,
};

UserItem.defaultProps = {
};

export default UserItem;
