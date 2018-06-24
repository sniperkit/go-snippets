import React from 'react';
import { Link } from 'react-router';

const UncrawlableItem = ({ data }) => {
  const uncrawlable = data;
  return (
    <div className="uncrawlable item col-md-12">
      <Link className="blue" to={`/uncrawlables/${uncrawlable.id}`}>{uncrawlable.url}</Link>
    </div>
  );
};

UncrawlableItem.propTypes = {
  data: React.PropTypes.object.isRequired,
  // onSelect: React.PropTypes.func.isRequired,
};

UncrawlableItem.defaultProps = {
};

export default UncrawlableItem;
