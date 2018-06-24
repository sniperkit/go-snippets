import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import PageHeader from '../PageHeader';

const Monitoring = ({ label, title, description }) => {
  return (
    <div id="monitoring" className="activity page">
      <PageHeader label="activity" title="Monitoring" />
      <div className="container">
      </div>
    </div>
  );
};

Monitoring.propTypes = {
};

Monitoring.defaultProps = {
};

export default Monitoring;
