import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import analytics from '../analytics';

import Consensus from '../components/Consensus';
import HistoryItem from '../components/item/HistoryItem';
import List from '../components/List';
// import Metadata from '../components/Metadata';
import StatsBar from '../components/StatsBar';
import TabBar from '../components/TabBar';
import UrlItem from '../components/item/UrlItem';
import MetadataEditor from './MetadataEditor';

import { loadConsensus } from '../actions/consensus';
import { loadContentMetadata, loadContentUrls } from '../actions/content';
import { selectConsensus } from '../selectors/consensus';
import { selectContentUrls, selectContentMetadata } from '../selectors/content';
import { selectUrlByHash } from '../selectors/url';
import { contentStats } from '../selectors/content';

class Content extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      tab: "metadata",
    };

    [
      "handleSetTab",
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('content');
    this.props.loadContentUrls(this.props.hash);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.hash != this.props.hash) {
      this.props.loadContentUrls(nextProps.hash);
    }
  }

  handleSetTab(tab) {
    const { hash } = this.props;
    switch (tab) {
      case "consensus":
        this.props.loadConsensus(hash);
        break;
      case "metadata":
        this.props.loadContentMetadata(hash);
        break;
      case "urls":
        this.props.loadContentUrls(hash);
        break;
      case "history":
        break;
      default:
        break;
    }
    this.setState({ tab });
  }

  renderCurrentTab() {
    const { tab } = this.state;
    switch (tab) {
      case "consensus":
        return <Consensus data={this.props.consensus ? this.props.consensus.data : {}} />;
      case "metadata":
        return <MetadataEditor subjectHash={this.props.hash} />;
      case "urls":
        return <List data={this.props.urls} component={UrlItem} />;
      case "history":
        return <List data={this.props.history} component={HistoryItem} />;
      default:
        return undefined;
    }
  }

  render() {
    const { hash, url } = this.props;
    const { tab } = this.state;

    return (
      <div id="content" className="content page">
        <header className="colorized">
          <div className="container">
            <div className="row">
              <div className="col-md-12">
                <hr className="" />
                <label className="label">Content</label>
                <h3 className="title">{url && url.fileName ? url.fileName : "unnamed content"}</h3>
                {url && <a className="btn" href={url.contentUrl} download={url.fileName}>Download File</a> }
                <hr />
              </div>
            </div>
            <StatsBar stats={contentStats(url)} />
          </div>
        </header>
        <div className="container">
          <TabBar value={tab} tabs={["metadata", "urls", "history", "consensus"]} onChange={this.handleSetTab} color="yellow" />
          <div className="row">
            {this.renderCurrentTab()}
          </div>
        </div>
      </div>
    );
  }
}

Content.propTypes = {
  hash: PropTypes.string.isRequired,

  consensus: PropTypes.object,
  // metadata: PropTypes.array,
  urls: PropTypes.array,
  history: PropTypes.array,

  loadContentUrls: PropTypes.func.isRequired,
  loadContentMetadata: PropTypes.func.isRequired,
  loadConsensus: PropTypes.func.isRequired,
};

Content.defaultProps = {
};

function mapStateToProps(state, ownProps) {
  const hash = ownProps.params.hash;

  return Object.assign({
    hash,
    url: selectUrlByHash(state, hash),
    urls: selectContentUrls(state, hash),
    consensus: selectConsensus(state, hash),
    metadata: selectContentMetadata(state, hash),
  }, ownProps);
}

export default connect(mapStateToProps, {
  loadConsensus,
  loadContentMetadata,
  loadContentUrls,
})(Content);
