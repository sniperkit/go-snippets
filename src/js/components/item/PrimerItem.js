import React from 'react';
import { Link } from 'react-router';

import ProgressBar from '../ProgressBar';



const PrimerItem = ({ data }) => {
  const primer = data;

  return (
    <div className="primer item col-md-4">
      <div className="card">
        <Link to={`/primers/${primer.id}`}>
          <header className="colorized">
            <h3>{primer.shortTitle}</h3>
          </header>
        </Link>
        <div className="info">
          <i>{primer.title}</i>
          {primer.stats &&
            <div className="stats">
              <b className="left">{primer.stats.urlCount} urls</b>
              <p className="right">{primer.stats.contentMetadataCount} / {primer.stats.contentUrlCount} completed</p>
            </div>
          }
        </div>
      </div>
    </div>
  );
};

PrimerItem.propTypes = {
  data: React.PropTypes.object.isRequired,
  // onSelect: React.PropTypes.func.isRequired,
};

PrimerItem.defaultProps = {
};

export default PrimerItem;
