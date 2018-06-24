import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import analytics from '../analytics';
import { selectSessionUser } from '../selectors/session';
import { selectCollectionItems } from '../selectors/collections';

import { 
  loadCollectionItems,
  saveCollectionItems,
  deleteCollectionItems,
} from '../actions/collections';

import Spinner from '../components/Spinner';
import CollectionItemItem from '../components/item/CollectionItemItem';

class CollectionItems extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
      allSelected: false,
      selected: [],
      editing: [],
      created: [],
    };

    [
      "handleAddItem",
      "handleDeleteSelectedItems",
      "handleNewItemChange",
      "handleRemoveNewItem",
      "handleSaveNewItems",
      "handleToggleAll",
      "handleToggleSelection",
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('collections');
    this.props.loadCollectionItems(this.props.id);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.items && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  handleToggleAll() {
    if (this.state.allSelected) {
      this.setState({
        allSelected: false,
        selected: [],
      });
    } else {
      this.setState({
        allSelected: true,
        selected: this.props.items.map((m, i) => i),
      });
    }
  }

  handleToggleSelection(index) {
    this.setState({
      selected: this.state.selected.concat([index]),
    })
  }

  handleDeleteSelectedItems() {
    if (confirm(`Are you sure you want to remove ${this.state.selected.length} ${(this.state.selected.length == 1) ? 'item' : 'items'} from this collection?`)) {
      const items = this.state.selected.map((i) => this.props.items[i])
      this.props.deleteCollectionItems(this.props.id, items, () => {
      });
    }
  }

  handleAddItem() {
    this.setState({
      created: this.state.created.concat([{ url: "", description: "" }]),
    });
  }

  handleRemoveNewItem(index, item) {
    const created = this.state.created.slice();
    created.splice(index, 1);
    this.setState({ created: created });
  }

  handleNewItemChange(index, item) {
    const created = this.state.created.slice();
    created.splice(index, 1, item);
    this.setState({ created });
  }

  handleSaveNewItems() {
    this.props.saveCollectionItems(this.props.id, this.state.created, () => {
      // clear created stuff
      this.setState({ created: [] });
    });
  }

  render() {
    const { loading, allSelected, selected, editing, created } = this.state;
    const { items, editable } = this.props;

    if (loading) {
      return <Spinner />;
    }

    return (
      <div className="container">
        <div className="row">
          <div className="col-md-12">
            <p>{ selected.length != 0 && <a onClick={this.handleDeleteSelectedItems}>Delete Items</a>} <a onClick={this.handleAddItem}>Add Item</a> &nbsp;</p>
          </div>
        </div>
        <div className="row">
          <div className="col-md-12">
            <table className="collection items table">
              <thead>
                <tr>
                  {editable && <th><input type="checkbox" onChange={this.handleToggleAll} value={allSelected} /></th>}
                  <th>hash</th>
                  <th>url</th>
                  <th>description</th>
                </tr>
              </thead>
              <tbody>
                {(created.length > 0) &&
                  created.map((item, i) => {
                    return (<CollectionItemItem 
                      key={`created-${i}`} 
                      index={i} 
                      editing={true} 
                      editable={editable} 
                      data={item} 
                      onChange={this.handleNewItemChange}
                      onDelete={this.handleRemoveNewItem} 
                    />);
                  })}
                {(created.length > 0) &&
                  <tr>
                    <td></td>
                    <td></td>
                    <td></td>
                    <td><button className="btn right btn-primary" onClick={this.handleSaveNewItems}>{(created.length == 1) ? "Add Url" : "Add Urls" }</button></td>
                  </tr>
                }
                {items.map((item,i) => {
                  return (<CollectionItemItem 
                    key={item.id}
                    editable={editable}
                    data={item}
                    index={i}
                    onToggleSelect={this.handleToggleSelection}
                    checked={selected.includes(i)}
                  />);
                })}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    );
  }
}

CollectionItems.propTypes = {
  id : PropTypes.string.isRequired,
  editable: PropTypes.bool,
  loadCollectionItems: PropTypes.func.isRequired,
  deleteCollectionItems: PropTypes.func.isRequired,
  saveCollectionItems: PropTypes.func.isRequired,
};

function mapStateToProps(state, ownProps) {
  return {
    id : ownProps.id,
    editable: ownProps.editable,
    user: selectSessionUser(state),
    items: selectCollectionItems(state, ownProps.id),
  };
}

export default connect(mapStateToProps, {
  loadCollectionItems,
  deleteCollectionItems,
  saveCollectionItems,
})(CollectionItems);
