/* -------------- STYLES ----------------------*/
const inputStyle = {
  width: '100%',
  'border' : '0px solid',
  'border-bottom': '1px dashed black'
}

/* ------------- UTILS ------------------------ */

function capitalize(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

function validText(text) {
  return text != null
    && text != ''
    && text.length > 2
    ;
}

function get_entries(data)
{
  let id = 0;
  return _.map(_.keys(data), url => {
    let entry = data[url];
    return {
      id: ++id,
      url: url,
      name: capitalize(entry['name']),
      desc: capitalize(entry['desc']),
      tags: entry['tags'],
    }
  });
}

function get_tags(entries)
{
    var all_tags = [];

    _.each(entries, entry => {
      let tags = entry["tags"];

      _.each(tags, tag => {
        if (_.findIndex(all_tags, e => e.name == tag) == -1) {
          all_tags.push({'name': tag, 'count': 0});
        }

        const pos = _.findIndex(all_tags, e => e.name == tag);

        all_tags[pos]['count']++;
      });

    });

    all_tags = _.sortBy(all_tags,'count');
    all_tags.reverse();

    return all_tags;
}

/* --------------- COMPONENTS -----------------*/
class Tag extends React.Component
{
  constructor(props) {
      super(props);

      this.handleClick = this.handleClick.bind(this);
  }

  handleClick(event) {
    this.props.onTagClick(event.target.id)
  }

  render() {
    let style={
      'font-size' : this.props.size + 'px',
      'padding-right': '10px',
      'display' : 'inline-block'
    }
    const name = this.props.name;
    return (
      <span style={style}><a href="#" id={name} onClick={this.handleClick}>{name}</a></span>
    );
  }
}

class TagCloud extends React.Component
{
  constructor(props) {
      super(props);
  }

  render() {

    if (this.props.tags != null && this.props.tags.length > 0 && this.props.show) {
      let minFont = 15;
      let maxFont = 40;
      let minCount = this.props.tags[this.props.tags.length-1]['count'];
      let maxCount = this.props.tags[0]['count'];
      let maxTags = this.props.tags.length; // 30

      let selTags = _.map(this.props.tags.slice(0, maxTags), e => {
        const v = (maxFont - minFont) * ( ( e['count'] - minCount ) / ( maxCount - minCount ) ) + minFont;
        return {
          'name': e['name'],
          'size': v
        }
      });

      return (
        <div>
          {selTags.map(tag =>
            <Tag name={tag['name']} size={tag['size']} onTagClick={this.props.onTagClick} />
          )}
        </div>
      );
    } else {
      return (
        <div></div>
      );
    }
  }
}

TagCloud.defaultProps = {
  'show': true,
}

class InputSearch extends React.Component
{
  render() {
    return (
      <div>
      <input type="search" style={inputStyle} placeholder="Search" onChange={this.props.handleChange} value={this.props.text} />
      </div>
    );
  }
}

InputSearch.defaultProps = {
  'text': '',
}


class ResultRow extends React.Component
{
    render() {
      const entry = this.props.entry;
      return (
        <div>
          <li>
            <a href={entry.url} target="_blank">
              <b>{entry.name}</b>: {entry.desc}
            </a>
          </li>
        </div>
      );
    }
}


class ResultLinks extends React.Component
{
  render() {
    const searchTerm = this.props.searchTerm;
    const entries = this.props.entries
    if (validText(searchTerm)) {
      const filtered = _.filter(entries, e => {
        return e['name'].toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1
          && e['desc'].toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1
      });

      let results = _.filter(entries, e => {
        return e['name'].toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1
          || e['desc'].toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1
          || _.contains(e['tags'], searchTerm.toLowerCase())
      });

      return (
        <div>
          <span>{results.length} results [<a href="#" onClick={this.props.clickReset}>Reset</a>]</span>
          <br/>
          <ul>
          {_.map(results, e =>
            <ResultRow entry={e} />
          )}
          </ul>
        </div>
      );
    }
    else {
      return (
        <div></div>
      );
    }
  }
}

class App extends React.Component
{
  constructor(props) {
      super(props);

      this.state = {
        entries: [],
        tags: [],
        textSearch: ''
      };

      this.handleChange = this.handleChange.bind(this);
      this.handleTagClick = this.handleTagClick.bind(this);
      this.handleResetClick = this.handleResetClick.bind(this);
  }

  componentDidMount() {
    axios.get("https://jepemo.github.io/links/db.json")
      .then(res => {
        let data = res.data;
        let entries = get_entries(data);
        let tags = get_tags(entries);

        this.setState({
          entries: entries,
          tags: tags
        });
      });
  }

  handleChange(event) {
    this.setState({textSearch: event.target.value})
  }

  handleTagClick(tagName) {
    this.setState({textSearch: tagName})
  }

  handleResetClick(event) {
    this.setState({ textSearch: '' });
  }

  render() {
    const showText = validText(this.state.textSearch);

    return(
    <div>
      <InputSearch handleChange={this.handleChange} text={this.state.textSearch} />
      <br/><br/>
      <TagCloud tags={this.state.tags} show={!showText} onTagClick={this.handleTagClick} />
      <ResultLinks searchTerm={this.state.textSearch} entries={this.state.entries} clickReset={this.handleResetClick} />
    </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('root'));
