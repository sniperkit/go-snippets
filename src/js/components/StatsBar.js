import React, { PropTypes } from 'react';

const StatsBar = ({ stats }) => {
  return (
    <div className="stats bar row">
      {Object.keys(stats).map((key, i) => {
        return (
          <span className="stat col-md-3" key={i}>
            <label>{key}</label>
            <h6>{stats[key]}&nbsp;</h6>
          </span>
        );
      })}
    </div>
  );
};

StatsBar.propTypes = {
  stats: PropTypes.object.isRequired,
};

export default StatsBar;
