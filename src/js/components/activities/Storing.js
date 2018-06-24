import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import PageHeader from '../PageHeader';

const Storing = ({ label, title, description }) => {
  return (
    <div id="storing" className="activity page">
      <PageHeader label="activity" title="Storing" />
      <div className="container">
      </div>
    </div>
  );
};

Storing.propTypes = {
};

Storing.defaultProps = {
};

export default Storing;
