import React from 'react';
import { Link } from 'react-router';

import ProgressBar from '../ProgressBar';

const SourceItem = ({ data }) => {
  const source = data;

  return (
    <div className="source item col-md-3">
      <Link to={`/sources/${source.id}`}>
        <h4 className="title">{source.title}</h4>
        <div className="info">
          <ProgressBar size="micro" color="blue" total={source.stats.contentUrlCount} progress={source.stats.contentMetadataCount} />
          <i>{source.stats.contentMetadataCount}/{source.stats.contentUrlCount}</i>
        </div>
      </Link>
    </div>
  );
};

SourceItem.propTypes = {
  data: React.PropTypes.object.isRequired,
  // onSelect: React.PropTypes.func.isRequired,
};

SourceItem.defaultProps = {
};

export default SourceItem;
