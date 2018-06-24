import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import PageHeader from '../PageHeader';

const Analyzing = ({ label, title, description }) => {
  return (
    <div id="analyzing" className="activity page">
      <PageHeader label="activity" title="Analyzing" />
      <div className="container">
      </div>
    </div>
  );
};

Analyzing.propTypes = {
};

Analyzing.defaultProps = {
};

export default Analyzing;
