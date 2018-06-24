import React, { PropTypes } from 'react';

const List = (props) => {
  // This strange props destructuring is so props.component
  // can be referenced as a jsx component below
  const { data, onSelectItem, className } = props;
  const selectFunc = (fn, d, i) => {
    return () => {
      if (fn) {
        fn(i, d);
      }
    };
  };

  return (
    <div className={className}>
      {data.map((item, i) => <props.component {...props} data={item} key={i} index={i} onSelect={selectFunc(onSelectItem, item, i)} />)}
    </div>
  );
};

List.propTypes = {
  data: PropTypes.array,
  // eslint-disable-next-line react/no-unused-prop-types
  component: PropTypes.func.isRequired,
  onSelectItem: PropTypes.func,
};

List.defaultProps = {
  data: [],
  className: "list",
};

export default List;
