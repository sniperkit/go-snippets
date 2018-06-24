import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link, browserHistory } from 'react-router';

import analytics from '../analytics';
import { selectCollection, selectLocalCollection } from '../selectors/collections';
import { selectAvailableUsers } from '../selectors/session';
import { selectCollectionActiveTasks } from '../selectors/tasks';

import { 
  loadCollection,
  newCollection,
  editCollection,
  updateCollection,
  saveCollection,
  deleteCollection,
  cancelCollectionEdit,
  loadCollectionItems,
  saveCollectionItems,
  deleteCollectionItems,
} from '../actions/collections';
import { archiveCollection } from '../actions/tasks';
import { selectDefaultKeyId } from '../selectors/keys';

import Spinner from '../components/Spinner';
import List from '../components/List';
import ContentItem from '../components/item/ContentItem';
import TaskItem from '../components/item/TaskItem';
import CollectionForm from '../components/form/CollectionForm';

import CollectionItems from './CollectionItems';

class Collection extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: !(props.id == "new"),
    };

    [
      "handleChange",
      "handleCancel",
      "handleSave",
      "handleEdit",
      "handleDelete",
      "handleArchiveCollection",
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('collection');
    if (this.props.id == "new") {
      this.props.newCollection(this.props.sessionKeyId);
    } else {
      this.props.loadCollection(this.props.id);
      this.props.loadCollectionItems(this.props.id, 1, 100);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.collection && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  handleChange(name, value) {
    this.props.updateCollection(value);
  }
  handleCancel() {
    this.props.cancelCollectionEdit(this.props.local);
  }
  handleEdit() {
    this.props.editCollection(this.props.collection);
  }
  handleSave() {
    if (this.props.local.id) {
      this.props.saveCollectionItems(this.props.local.id,this.props.local.contents);
    }
    this.props.saveCollection(this.props.local, (collection) => {
      browserHistory.push(`/collections/${collection.id}`);
    });
  }
  handleDelete() {
    if (confirm("are you sure you want to delete this collection?")) {
      this.props.deleteCollection(this.props.collection, (collection) => {
        browserHistory.push(`/collections`);
      });
    }
  }
  handleArchiveCollection() {
    this.props.archiveCollection(this.props.collection);
  }

  render() {
    const { loading } = this.state;
    const { collection, items, local, users, sessionKeyId, id, activeTasks } = this.props;

    if (loading) {
      return <Spinner />;
    } else if (!collection && !local) {
      return null;
    } else if (local) {
      return (
        <CollectionForm
          data={local}
          users={users}
          onChange={this.handleChange}
          onCancel={this.handleCancel}
          onSubmit={this.handleSave}
        />
      );
    }

    return (
      <div id="collection" className="collection page">
        <header className="collection colorized">
          <div className="container">
            <div className="row">
              <div className="col-md-12">
                <hr className="green" />
                <a className="right" onClick={this.handleDelete}>&nbsp; Delete</a>
                <a className="right" onClick={this.handleEdit}>&nbsp; Edit</a>
                {sessionKeyId && <a className="right red" onClick={this.handleArchiveCollection}>&nbsp; Archive</a>}
                <label className="label">Collection</label>
                <h1>{collection.title}</h1>
                <p>{collection.description}</p>
              </div>
            </div>
          </div>
        </header>
        <div className="container">
          <div className="row">
            <div className="col-md-12">
              {activeTasks.length > 0 && <label className="label">tasks</label>}
            </div>
            <List data={activeTasks} component={TaskItem} />
          </div>
        </div>
        <CollectionItems id={id} editable={true} />
      </div>
    );
  }
}

Collection.propTypes = {
  id: PropTypes.string.isRequired,
  sessionKeyId: PropTypes.string,

  collection: PropTypes.object,
  items : PropTypes.array,
  local: PropTypes.object,
  loadCollection: PropTypes.func.isRequired,
  // editCollection: PropTypes.func.isRequired,
};

function mapStateToProps(state, ownProps) {
  const { id } = ownProps.params;
  const sessionKeyId = selectDefaultKeyId(state);

  return {
    id,
    sessionKeyId,
    users : selectAvailableUsers(state),
    local: selectLocalCollection(state, id),
    collection: selectCollection(state, id),
    activeTasks : selectCollectionActiveTasks(state, id),
  };
}

export default connect(mapStateToProps, {
  loadCollection,
  newCollection,
  editCollection,
  updateCollection,
  saveCollection,
  deleteCollection,
  cancelCollectionEdit,
  archiveCollection,
  loadCollectionItems,
  saveCollectionItems,
})(Collection);
