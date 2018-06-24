import * as React from 'react';
import 'semantic-ui-css/semantic.min.css';
import './stylesheets/theme.css';
import './stylesheets/App.css';
import {
    ActionBar,
    ActionBarRow,
    HitItemProps,
    Hits,
    HitsStats,
    Layout,
    LayoutBody,
    LayoutResults,
    NoHits,
    SearchBox,
    SearchkitManager,
    SearchkitProvider,
    TopBar,
} from 'searchkit';
import { get } from 'lodash';

// Set ES url - use a protected URL that only allows read actions.
const elasticsearchHost = 'http://localhost:9200';
const searchkit = new SearchkitManager(elasticsearchHost);

const services = new Map([
    ['confluence', '/images/confluence.png'],
    ['gdrive', '/images/gdrive.png'],
    ['github', '/images/github.svg'],
]);

// Note: This is a very hacky way to highlight search string matches due to weird
// styling interactions between Searchkit's theming and Semantic UI.
const highlightMatches = (s: string) => s.replace('<em>', '<mark>').replace('</em>', '</mark>');

const HitItem = (props: HitItemProps) => {
    const { result } = props;
    const highlightedTitle = get(result, 'highlight.title', false);
    const title = highlightedTitle ?
        highlightMatches(highlightedTitle[0]) :
        result._source.title;
    const highlightedDescription = get(result, 'highlight.description', false);
    const description = highlightedDescription ?
        highlightMatches(highlightedDescription[0]) :
        result._source.description;
    const highlightedContent = get(result, 'highlight.content', false);
    const content = highlightedContent && highlightedContent[0] !== highlightedDescription[0] ?
        highlightMatches(highlightedContent[0]) :
        '';

    return (
        <div className="ui link items">
            <a className="item" href={result._source.url}>
                <div className="ui tiny image">
                    <img src={services.get(result._source.service)}/>
                </div>
                <div className="content">
                    <div
                        className="header"
                        dangerouslySetInnerHTML={{__html: title}}
                    />
                    <div
                        className="description"
                        dangerouslySetInnerHTML={{__html: description}}
                    />
                    <div
                        className="meta"
                        dangerouslySetInnerHTML={{__html: content}}
                    />
                    <div className="extra">{result._source.url}</div>
                </div>
            </a>
        </div>);
};

const App: React.SFC<{}> = () => (
    <SearchkitProvider searchkit={searchkit}>
        <Layout>
            <TopBar>
                <SearchBox
                    autofocus={true}
                    searchOnChange={false}
                    queryOptions={{analyzer: 'standard'}}
                    queryFields={['title', 'description', 'content']}
                    prefixQueryFields={['title^10', 'description^2', 'content']}
                />
            </TopBar>
            <LayoutBody>
                <LayoutResults>
                    <ActionBar>
                        <ActionBarRow>
                            <HitsStats/>
                        </ActionBarRow>
                    </ActionBar>
                    {/*TODO: Add pagination.*/}
                    <Hits
                        hitsPerPage={25}
                        highlightFields={['title', 'description', 'content']}
                        sourceFilter={['title', 'description', 'url', 'service']}
                        itemComponent={HitItem}
                    />
                    <NoHits/>
                </LayoutResults>
            </LayoutBody>
        </Layout>
    </SearchkitProvider>
);

export default App;
