import React from 'react';
import { Link } from 'react-router';

import { fileSizeString, fileTypeString } from '../../selectors/format';

const ContentItem = ({ data }) => {
  const content = data;
  return (
    <div className="content item col-md-12">
      <Link to={`/content/${content.hash}`}><h5 className="title">{content.fileName || "unnamed content"}</h5></Link>
      <p>{fileSizeString(content.contentLength)} | {fileTypeString(content.contentSniff)}</p>
    </div>
  );
};

ContentItem.propTypes = {
  data: React.PropTypes.object.isRequired,
  // onSelect: React.PropTypes.func.isRequired,
};

ContentItem.defaultProps = {
};

export default ContentItem;
