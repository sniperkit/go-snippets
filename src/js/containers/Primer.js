import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import analytics from '../analytics';
import { selectPrimer } from '../selectors/primers';
import { loadPrimer } from '../actions/primers';

import Spinner from '../components/Spinner';
import List from '../components/List';
import SourceItem from '../components/item/SourceItem';
import PrimerItem from '../components/item/PrimerItem';
import StatsBar from '../components/StatsBar';
import ProgressBar from '../components/ProgressBar';

class Primer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
    };

    [].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('primer');
    this.props.loadPrimer(this.props.id);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.primer && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  renderSubPrimers() {
    const { primer } = this.props;
    if (!primer.subPrimers) {
      return undefined;
    }

    return (
      <div className="primer container">
        <div className="row">
          <div className="col-md-12">
            <h3 className="title">Sub-Primers:</h3>
            <br />
          </div>
        </div>
        <div className="row">
          <List data={primer.subPrimers} component={PrimerItem} />
        </div>
      </div>
    );
  }


  renderSources() {
    const { primer } = this.props;
    if (!primer.sources) {
      return undefined;
    }

    return (
      <div className="source container">
        <div className="row">
          <div className="col-md-12">
            <h3>Sources:</h3>
            <br />
          </div>
        </div>
        <div className="row">
          <List data={primer.sources} component={SourceItem} />
        </div>
      </div>
    );
  }

  render() {
    const { loading } = this.state;
    const { primer } = this.props;

    if (loading) {
      return <Spinner />;
    }

    return (
      <div id="primers" className="primer page">
        <header className="colorized">
          <div className="container">
            <div className="row">
              <div className="col-md-12">
                <hr/>
                <label className="label">Primer</label>
                <h1 className="title">{primer.title}</h1>
              </div>
            </div>
            <div className="row">
              <div className="col-md-12 col-lg-8">
                <p>{primer.description}</p>
              </div>
            </div>
            <div className="row">
              <div className="col-md-12">
                <ProgressBar 
                  total={primer.stats.contentUrlCount}
                  progress={primer.stats.contentMetadataCount}
                />
                <StatsBar
                  stats={{
                    urls: primer.stats.urlCount,
                    content: primer.stats.contentUrlCount,
                    documented: primer.stats.contentMetadataCount,
                  }}
                />
              </div>
            </div>
          </div>
        </header>
        {this.renderSubPrimers()}
        {this.renderSources()}
      </div>
    );
  }
}

Primer.propTypes = {
  id: PropTypes.string.isRequired,
  primer: PropTypes.object,
  loadPrimer: PropTypes.func.isRequired,
};

function mapStateToProps(state, ownProps) {
  return {
    id: ownProps.params.id,
    primer: selectPrimer(state, ownProps.params.id),
  };
}

export default connect(mapStateToProps, {
  loadPrimer,
})(Primer);
