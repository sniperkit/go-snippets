import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import analytics from '../analytics';
import { selectTasks } from '../selectors/tasks';
import { loadTasks, enqueueTask } from '../actions/tasks';

import List from '../components/List';
import TaskForm from '../components/form/TaskForm';
import TaskItem from '../components/item/TaskItem';
import Spinner from '../components/Spinner';

class Tasks extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
    };

    [].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('tasks');
    this.props.loadTasks(1, 25, true);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.tasks && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  render() {
    const { loading } = this.state;
    const { tasks } = this.props;

    if (loading) {
      return <Spinner />;
    }

    return (
      <div id="tasks" className="task page">
         <header className="container">
          <div className="row">
            <div className="col-md-12">
              <hr />
              <h1>Tasks:</h1>
            </div>
            <div className="col-md-6">
              <p>Hey look a queue of tasks</p>
            </div>
            <div className="col-md-12">
              <hr className="" />
            </div>
          </div>
        </header>
        <div className="container">
          <div className="row">
            <br />
            <List data={tasks} component={TaskItem} />
          </div>
        </div>
        <div className="container">
          <TaskForm onSubmit={this.props.enqueueTask} />
        </div>
      </div>
    );
  }
}

Tasks.propTypes = {
  // user: PropTypes.object,
  tasks: PropTypes.array,
  loadTasks: PropTypes.func.isRequired,
  enqueueTask: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    tasks: selectTasks(state),
  };
}

export default connect(mapStateToProps, {
  loadTasks,
  enqueueTask,
})(Tasks);
