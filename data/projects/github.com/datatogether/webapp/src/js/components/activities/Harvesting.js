import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import PageHeader from '../PageHeader';

const Harvesting = ({ label, title, description }) => {
  return (
    <div id="harvesting" className="activity page">
      <PageHeader label="activity" title="Harvesting" />
      <div className="container">
      </div>
    </div>
  );
};

Harvesting.propTypes = {
};

Harvesting.defaultProps = {
};

export default Harvesting;
