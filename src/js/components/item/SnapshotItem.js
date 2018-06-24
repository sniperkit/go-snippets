import React from 'react';
import { Link } from 'react-router';

const SnapshotItem = ({ data }) => {
  return (
    <div className="search result item col-md-12">
      <p>{data.hash ? <Link className="hash" to={`/content/${data.hash}`}>{data.hash}</Link> : undefined }</p>
    </div>
  );
};

SnapshotItem.propTypes = {
  data: React.PropTypes.object.isRequired,
  // onSelect: React.PropTypes.func.isRequired,
};

SnapshotItem.defaultProps = {
};

export default SnapshotItem;
