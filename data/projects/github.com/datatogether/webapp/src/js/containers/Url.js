import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';

import analytics from '../analytics';
import { selectUrl, selectOutboundLinks, selectInboundLinks, urlStats, containsContent } from '../selectors/url';
import { loadUrl, loadOutboundLinks, loadInboundLinks, archiveUrl } from '../actions/url';

import List from '../components/List';
import NotFound from '../components/NotFound';
import SnapshotItem from '../components/item/SnapshotItem';
import Spinner from '../components/Spinner';
import StatsBar from '../components/StatsBar';
import UrlItem from '../components/item/UrlItem';
import TabBar from '../components/TabBar';

function defaultTab(props) {
  let tab = 'outbound links';
  if (props.url && containsContent(props.url)) {
    tab = 'content';
  }
  return tab;
}

class Url extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: (props.url == undefined),
      tab: defaultTab(props),
    };

    [
      'handleArchive',
      'handleUncrawlableClick',
      'handleChangeTab',
    ].forEach((m) => { this[m] = this[m].bind(this); });
  }

  componentWillMount() {
    analytics.page('url');
    this.props.loadUrl(this.props.urlParam);
    this.props.loadInboundLinks(this.props.urlParam);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.urlParam != this.props.urlParam) {
      this.props.loadUrl(nextProps.urlParam);
      this.props.loadOutboundLinks(this.props.urlParam);
      this.setState({ loading: true, tab: defaultTab(nextProps) });
    } if (nextProps.url && !this.props.url) {
      this.setState({ loading: false, tab: defaultTab(nextProps) });
    }
  }

  handleArchive() {
    this.props.archiveUrl(this.props.url.url);
  }

  handleUncrawlableClick() {
    alert("marking links as uncrawlable coming soon!");
  }

  handleChangeTab(tab) {
    const { urlParam } = this.props;
    switch (tab) {
      case "outbound links":
        this.props.loadOutboundLinks(urlParam);
        break;
      case "inbound links":
        this.props.loadInboundLinks(urlParam);
        break;
      case "snapshots":
        break;
      default:
        break;
    }
    this.setState({ tab });
  }

  renderCurrentTab() {
    const { tab } = this.state;
    const { url } = this.props;

    switch (tab) {
      case "outbound links":
        return <List data={this.props.outboundLinks} component={UrlItem} />;
      case "inbound links":
        return <List data={this.props.inboundLinks} component={UrlItem} />;
      case "snapshots":
        return <List data={this.props.snapshots} component={SnapshotItem} />;
      default:
        return undefined;
    }
  }

  renderContent() {
    const { url } = this.props;
    const { tab } = this.state;

    let tabs = ["outbound links", "inbound links", "snapshots"];

    if (containsContent(url)) {
      tabs = ["inbound links"];
    }

    return (
      <div className="container">
        {containsContent(url) &&
          <div className="row">
            <div className="col-md-12">
              <br />
              <label>Content</label><br />
              <b><Link className="content" to={`/content/${url.hash}`}>{url.fileName || url.hash}</Link></b>
            </div>
          </div>
        }
        <TabBar value={tab} tabs={tabs} onChange={this.handleChangeTab} color="blue" />
        <div className="row">
          {this.renderCurrentTab()}
        </div>
      </div>
    );
  }

  render() {
    const { url } = this.props;
    const { loading } = this.state;

    if (loading) {
      return <Spinner />;
    } else if (!url) {
      return <NotFound />;
    }

    return (
      <div id="url" className="url page">
        <header className="colorized">
          <div className="container">
            <div className="row">
              <div className="col-md-12">
                <hr className="blue" />
                <a className="right" target="_blank" rel="noopener noreferrer" href={url.url}>link</a>
                <label className="label">URL</label>
                <h4 className="title">{url.title}</h4>
                <h5>{url.url}</h5>
                <br />
                {!url.hash && <button className="btn btn-primary" onClick={this.handleArchive}>Archive Url</button>}
                {url.hash && !containsContent(url) && <button className="btn" onClick={this.handleUncrawlableClick}>Missing Content</button>}
                <hr className="blue" />
              </div>
            </div>
            <StatsBar stats={urlStats(url)} />
          </div>
        </header>
        {this.renderContent()}
      </div>
    );
  }
}

Url.propTypes = {
  url: PropTypes.object,
  loadUrl: PropTypes.func.isRequired,
  outboundLinks: PropTypes.array,
  inboundLinks: PropTypes.array,
  snapshots: PropTypes.array,
};

function mapStateToProps(state, ownProps) {
  const urlParam = ownProps.location.query.url;

  return Object.assign({
    // user : selectLocalSessionUser(state),
    urlParam,
    url: selectUrl(state, urlParam),
    outboundLinks: selectOutboundLinks(state, urlParam),
    inboundLinks: selectInboundLinks(state, urlParam),
  }, ownProps);
}

export default connect(mapStateToProps, {
  loadUrl,
  loadOutboundLinks,
  loadInboundLinks,
  archiveUrl,
})(Url);
