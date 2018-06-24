import React, { PropTypes } from 'react';
import { Link } from 'react-router';

const PageHeader = ({ label, title, description }) => {
  return (
    <header className="colorized">
      <div className="container">
        <div className="row">
          <div className="col-md-12">
            <hr className="" />
            {label && <label className="label">{label}</label>}
            {title && <h3 className="title">{title}</h3>}
            {description && <p>{description}</p>}
            <hr />
          </div>
        </div>
      </div>
    </header>
  );
};

PageHeader.propTypes = {
  label: PropTypes.string,
  title: PropTypes.string,
  description: PropTypes.string,
};

PageHeader.defaultProps = {
};

export default PageHeader;
