import React from 'react';

const SessionRequired = () => {
  return (
    <div id="wrapper" className="sessionRequired">
      <div className="container">
        <div className="col-md-12">
          <h3 className="center-text">You must be logged-in to do this</h3>
        </div>
      </div>
    </div>
  );
};

export default SessionRequired;
