import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import analytics from '../analytics';
import { selectSource, selectSourceUndescribedUrls, selectSourceAttributedUrls } from '../selectors/sources';
import { loadSource, loadSourceUrls, loadSourceAttributedUrls } from '../actions/source';

import List from '../components/List';
import ProgressBar from '../components/ProgressBar';
import Spinner from '../components/Spinner';
import StatsBar from '../components/StatsBar';
import TabBar from '../components/TabBar';
import UrlItem from '../components/item/UrlItem';

class Source extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: (props.source == undefined),
      tab: 'unattributed content',
    };

    [
      "handleChangeTab",
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('source');
    this.props.loadSource(this.props.id);
    this.props.loadSourceUrls(this.props.id);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.source && this.state.loading) {
      this.setState({ loading: false });
    }
  }

  handleChangeTab(tab) {
    switch (tab) {
      case "unattributed content":
        this.props.loadSourceUrls(this.props.id);
        break;
      case "attributed content":
        this.props.loadSourceAttributedUrls(this.props.id);
        break;
      default:
        break;
    }
    this.setState({ tab });
  }

  renderCurrentTab() {
    const { tab } = this.state;
    const { urls, attributedUrls } = this.props;


    switch (tab) {
      case "unattributed content":
        return (
          <div className="row">
            <div className="col-md-12">
              <h4>Content Needing Metadata:</h4>
            </div>
            <List data={urls} component={UrlItem} />
          </div>);
      case "attributed content":
        return (
          <div className="row">
            <div className="col-md-12">
              <h4>Content With Metadata:</h4>
            </div>
            <List data={attributedUrls} component={UrlItem} />
          </div>);
      default:
        return undefined;
    }
  }

  render() {
    const { loading, tab } = this.state;
    const { source } = this.props;

    if (loading) {
      return <Spinner />;
    }

    return (
      <div id="source" className="source page">
        <header className="colorized">
          <div className="container">
            <div className="row">
              <div className="col-md-12">
                <hr />
                <label className="label">Source</label>
                <h1 className="title">{source.title}</h1>
                <p>{source.url}</p>
                <label className="label">Primer:</label>
                <Link className="primer" to={`/primers/${source.primer.id}`}>
                  <h5 className="title">{source.primer.title}</h5>
                </Link>
              </div>
            </div>
            <div className="row">
              <div className="col-md-12">
                <ProgressBar 
                  total={source.stats.contentUrlCount}
                  progress={source.stats.contentMetadataCount}
                />
                <StatsBar
                  stats={{
                    urls: source.stats.urlCount,
                    content: source.stats.contentUrlCount,
                    documented: source.stats.contentMetadataCount,
                  }}
                />
              </div>
            </div>
          </div>
        </header>
        <div className="container">
          <div className="row">
            <div className="col-md-12">
              <TabBar value={tab} tabs={["unattributed content", "attributed content"]} onChange={this.handleChangeTab} color="orange" />
            </div>
          </div>
          {this.renderCurrentTab()}
        </div>
      </div>
    );
  }
}

Source.propTypes = {
  id: PropTypes.string.isRequired,
  source: PropTypes.object,

  urls: PropTypes.array.isRequired,
  attributedUrls: PropTypes.array.isRequired,

  loadSource: PropTypes.func.isRequired,
  loadSourceUrls: PropTypes.func.isRequired,
};

function mapStateToProps(state, ownProps) {
  return {
    id: ownProps.params.id,
    source: selectSource(state, ownProps.params.id),

    urls: selectSourceUndescribedUrls(state, ownProps.params.id),
    attributedUrls: selectSourceAttributedUrls(state, ownProps.params.id),
  };
}

export default connect(mapStateToProps, {
  loadSource,
  loadSourceUrls,
  loadSourceAttributedUrls,
})(Source);
