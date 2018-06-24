import React, { PropTypes } from 'react';

const TabBar = ({ tabs, value, onChange, color }) => {
  const handleClick = (tab) => {
    return () => {
      onChange(tab);
    };
  };

  return (
    <div className="tab bar row">
      <div className="col-md-12">
        <hr className={color} />
        {tabs.map((tab, i) => {
          return (
            <a
              className={(value == tab) ? "current tab" : "tab"}
              onClick={handleClick(tab)}
              key={i}
            >{tab}</a>
          );
        })}
        <hr className={color} />
      </div>
    </div>
  );
};

TabBar.propTypes = {
  color: PropTypes.string,
  tabs: PropTypes.array.isRequired,
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
};

TabBar.defaultProps = {
  color: "",
};

export default TabBar;
