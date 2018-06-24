import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import PageHeader from '../PageHeader';

const Activities = ({ label, title, description }) => {
  return (
    <div id="activities" className="activity page">
      <PageHeader title="Activities" />
      <div className="container">
        <div className="row">
          <div className="col-md-12">
            <h2><Link className="red" to="/activities/harvesting">Harvesting</Link></h2>
            <br />
            <h2><Link className="red" to="/activities/monitoring">Monitoring</Link></h2>
            <br />
            <h2><Link className="red" to="/activities/storing">Storing</Link></h2>
            <br />
            <h2><Link className="red" to="/activities/analyzing">Analyzing</Link></h2>
            <br />
          </div>
        </div>
      </div>
    </div>
  );
};

Activities.propTypes = {
};

Activities.defaultProps = {
};

export default Activities;
