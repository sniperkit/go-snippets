import React from 'react';
import { Link } from 'react-router';

const HistoryItem = ({ data }) => {
  const uri = encodeURIComponent(data.url);
  return (
    <div className="user item col-xl-3 col-lg-4 col-md-6 col-sm-12">
      <Link to={`/url?url=${uri}`}>{data.url}</Link>
      {data.hash ? <Link to={`/content/${data.hash}`}>{data.hash}</Link> : undefined }
    </div>
  );
};

HistoryItem.propTypes = {
  data: React.PropTypes.object.isRequired,
  // onSelect: React.PropTypes.func.isRequired,
};

HistoryItem.defaultProps = {
};

export default HistoryItem;
