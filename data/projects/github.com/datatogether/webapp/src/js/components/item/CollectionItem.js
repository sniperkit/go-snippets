import React from 'react';
import { Link } from 'react-router';

const CollectionItem = ({ data }) => {
  const collection = data;
  return (
    <div className="collection item">
      <h4><Link className="green" to={`/collections/${collection.id}`}>{collection.title}</Link></h4>
      {collection.contents && <p>{collection.contents.length} {(collection.contents.length == 1) ? "url" : "urls"}</p>}
    </div>
  );
};

CollectionItem.propTypes = {
  data: React.PropTypes.object.isRequired,
  // onSelect: React.PropTypes.func.isRequired,
};

CollectionItem.defaultProps = {
};

export default CollectionItem;
