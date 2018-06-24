import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link, browserHistory } from 'react-router';
import { debounce } from 'lodash';

import analytics from '../analytics';

import { loadUrl } from '../actions/url';
import { setSearch, clearSearch } from '../actions/search';
import { archiveUrl } from '../actions/tasks';
import { selectDefaultKeyId } from '../selectors/keys';
import { selectUrlActiveTasks } from '../selectors/tasks';
import { selectUrl } from '../selectors/url';

import Spinner from '../components/Spinner';
import List from '../components/List';
import TaskItem from '../components/item/TaskItem';
import UrlItem from '../components/item/UrlItem';
import ValidInput from '../components/form/ValidInput';
import ValidTextarea from '../components/form/ValidTextarea';

class ArchiveUrl extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      usage: "",
      query: "",
    };

    [
      "handleChange",
      "handleArchiveUrl",
      "handleUrlChange"
    ].forEach((m) => { this[m] = this[m].bind(this); });

    // debounce url change func to avoid clobbering server
    this.handleUrlChange = debounce(this.handleUrlChange, 200);
  }

  componentWillMount() {
    analytics.page('archiveUrl');
  }

  componentWillReceiveProps(nextProps) {
  }

  handleChange(name, value) {
    this.setState({ [name]: value });
    if (name == "query") {
      this.handleUrlChange();
    }
  }

  handleSelectTask(index, data) {
    this.props.browserHistory.push(`/urls?url=${data.params.url}`);
  }

  handleUrlChange() {
    this.props.setSearch(this.state.query);
    this.props.loadUrl(this.state.query, true);
  }

  handleArchiveUrl() {
    this.props.archiveUrl(this.state.query);
  }

  render() {
    const { loading, query, usage } = this.state;
    const { activeTasks, url } = this.props;

    return (
      <div className="page">
        <div className="container">
          <div className="row">
            <div className="col-md-6 offset-md-3">
              <hr />
              <h4>Add a Dataset</h4>
              <ValidInput
                name="query"
                label="URL to archive"
                value={query}
                onChange={this.handleChange}
                placeholder="http://dataseturl.gov"
              />
              <ValidTextarea
                name="usage"
                label="How do you use this dataset?"
                value={usage}
                placeholder="I use this data for ..."
                onChange={this.handleChange}
              />
              <button className="btn btn-primary" onClick={this.handleArchiveUrl}>Add</button>
            </div>
          </div>
        </div>
        {url && <div className="container">
          <div className="row">
            <div className="col-md-6 offset-md-3">
              <hr />
              <label className="label">Found Urls:</label>
              <UrlItem data={url} />
            </div>
          </div>
        </div>}
        {activeTasks.length > 0 && <div className="container">
          <div className="row">
            <div className="col-md-6 offset-md-3">
              <hr />
              <label className="label">Archive Url Tasks</label>
              <List data={activeTasks} component={TaskItem} onSelect={this.handleSelectTask} />
            </div>
          </div>
        </div>}
      </div>
    );
  }
}

ArchiveUrl.propTypes = {
  sessionKeyId: PropTypes.string,
  archiveUrl: PropTypes.func,
};

function mapStateToProps(state, ownProps) {
  const sessionKeyId = selectDefaultKeyId(state);

  return {
    sessionKeyId,
    activeTasks: selectUrlActiveTasks(state),
    url: selectUrl(state, state.app.query),
  };
}

export default connect(mapStateToProps, {
  loadUrl,
  archiveUrl,
  setSearch,
  browserHistory,
})(ArchiveUrl);
