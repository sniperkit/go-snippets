import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';

import analytics from '../analytics';
import { selectUncrawlable, selectLocalUncrawlable } from '../selectors/uncrawlables';
import { loadUncrawlable, newUncrawlable, editUncrawlable } from '../actions/uncrawlables';
import { selectDefaultKeyId } from '../selectors/keys';

import Spinner from '../components/Spinner';
import List from '../components/List';
import ContentItem from '../components/item/ContentItem';
// import UncrawlableForm from '../components/form/UncrawlableForm';
// import UncrawlableView from '../components/Uncrawlable';

class Uncrawlable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: !(props.id == "new"),
    };

    [

    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('uncrawlable');
    (this.props.id == "new") ?
      this.props.newUncrawlable(this.props.sessionKeyId):
      this.props.loadUncrawlable(this.props.id);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.uncrawlable && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  handleChange(name, value) {

  }
  handleCancel() {
    
  }
  handleSave() {
    
  }

  render() {
    const { loading } = this.state;
    const { uncrawlable, local, sessionKeyId } = this.props;

    return (
      <div id="uncrawlable" className="page">
        <div className="container">
          <header className="row">
            <div className="col-md-12">
              <hr className="blue" />
              <label className="label">Uncrawlable</label>
              <h1 className="blue">{uncrawlable.url}</h1>
              <button className="btn btn-large bg-blue white">Link Github Repo</button>
            </div>
          </header>
          <div className="row">
            <div className="col-md-12 col-lg-8">
            </div>
          </div>
        </div>
      </div>
    );

    // if (loading) {
    //   return <Spinner />;
    // } else if (!uncrawlable && !local) {
    //   return null;
    // } else if (local) {
    //   return (
    //     <UncrawlableForm
    //       data={local}
    //       onChange={this.handleChange}
    //       onCancel={this.handleCancel}
    //       onSubmit={this.handleSave}
    //     />
    //   );
    // }

    // return (
    //   <UncrawlableView data={uncrawlable} />
    // );
  }
}

Uncrawlable.propTypes = {
  id: PropTypes.string.isRequired,
  sessionKeyId: PropTypes.string,

  uncrawlable: PropTypes.object,
  local: PropTypes.object,
  loadUncrawlable: PropTypes.func.isRequired,
  // editUncrawlable: PropTypes.func.isRequired,
};

function mapStateToProps(state, ownProps) {
  const { id } = ownProps.params;
  const sessionKeyId = selectDefaultKeyId(state);

  return {
    id,
    sessionKeyId,
    local: selectLocalUncrawlable(state, id),
    uncrawlable: selectUncrawlable(state, id),
  };
}

export default connect(mapStateToProps, {
  loadUncrawlable,
  newUncrawlable,
  editUncrawlable,
})(Uncrawlable);
