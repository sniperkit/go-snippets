import React from 'react';

import UrlItem from './UrlItem';

const SearchResultItem = ({ data }) => {
  return <UrlItem data={data} />;
};

SearchResultItem.propTypes = {
  data: React.PropTypes.object.isRequired,
  // onSelect: React.PropTypes.func.isRequired,
};

SearchResultItem.defaultProps = {
};

export default SearchResultItem;
