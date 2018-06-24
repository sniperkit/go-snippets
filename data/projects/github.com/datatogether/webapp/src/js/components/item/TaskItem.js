import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import ProgressBar from '../ProgressBar';

const TaskItem = ({ data, onSelect }) => {
  const task = data;

  return (
    <div className="task item">
      <Link to={task.progress && task.progress.dest}>
        <h6 className="title" onClick={onSelect}>{task.title}</h6>
        {(task.progress && !task.succeeded && !task.error) &&
        <div className="info">
          <ProgressBar size="micro" color="blue" progress={task.progress.percent * 100} />
          <i>{task.progress.status}</i>
        </div>}
        {(task.progress && task.succeeded && !task.error) && <div>done!</div>}
        {task.error && <div className="error">{task.error}</div>}
      </Link>
    </div>
  );
};

TaskItem.propTypes = {
  data: PropTypes.object.isRequired,
  onSelect: PropTypes.func.isRequired,
};

TaskItem.defaultProps = {
};

export default TaskItem;
