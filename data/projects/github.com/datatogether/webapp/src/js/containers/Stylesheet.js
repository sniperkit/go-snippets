import React from 'react';
import { connect } from 'react-redux';

import List from '../components/List';
import UserItem from '../components/item/UserItem';
import PrimerItem from '../components/item/PrimerItem';
import SourceItem from '../components/item/SourceItem';
import UrlItem from '../components/item/UrlItem';
import ContentItem from '../components/item/ContentItem';

const Stylesheet = (props) => {
  return (
    <div id="stylesheet" className="page">
      <div className="container">
        <div className="row">
          <header className="blue col-md-12">
            <hr className="blue" />
            <h1>Le Styles</h1>
          </header>
        </div>
        <div className="row">
          <div className="col-md-6">
            <h1>Heading One</h1>
            <h2>Heading Two</h2>
            <h3>Heading Three</h3>
            <h4>Heading Four</h4>
            <h5>Heading Five</h5>
            <h6>Heading Six</h6>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus a dui at enim sodales gravida.
            Donec imperdiet sagittis libero, a egestas nisl rhoncus nec. Proin massa elit, commodo a rutrum eget, condimentum ac lectus.
            Phasellus eget pulvinar tellus. Aliquam nec mi malesuada, feugiat augue molestie, viverra mi. Fusce non arcu venenatis, aliquam arcu eget, blandit risus. Ut id dui a turpis tristique luctus. Etiam dui risus, tincidunt eget elit id, lacinia molestie velit. Phasellus eleifend nunc vel vestibulum dapibus. Phasellus auctor et lectus sit amet fermentum. In eget ex ut dolor egestas tempor. Nam semper libero a elementum bibendum.
            Phasellus cursus, felis eu aliquet maximus, neque lorem egestas est, eu varius nulla sapien vitae tellus.</p>
          </div>
          <div className="col-md-6">
            <button className="btn btn-primary">Primary Button</button>
          </div>
        </div>
      </div>
      <div className="container">
        <div className="row">
          <header className="col-md-12">
            <hr className="red" />
            <h2>Palette</h2>
          </header>
        </div>
        <div className="row">
          <div className="col-md-2">
            <div className="swatch bg-green" style={{ width: 100, height: 100 }}></div>
            <label>Green</label>
          </div>
          <div className="col-md-2">
            <div className="swatch bg-blue" style={{ width: 100, height: 100 }}></div>
            <label>Blue</label>
          </div>
          <div className="col-md-2">
            <div className="swatch bg-purple" style={{ width: 100, height: 100 }}></div>
            <label>Purple</label>
          </div>
          <div className="col-md-2">
            <div className="swatch bg-red" style={{ width: 100, height: 100 }}></div>
            <label>Red</label>
          </div>
          <div className="col-md-2">
            <div className="swatch bg-yellow" style={{ width: 100, height: 100 }}></div>
            <label>Yellow</label>
          </div>
          <div className="col-md-2">
            <div className="swatch bg-orange" style={{ width: 100, height: 100 }}></div>
            <label>Orange</label>
          </div>
        </div>
      </div>
      <div className="container">
        <div className="row">
          <header className="col-md-12">
            <hr className="yellow" />
            <h2 className="yellow">Users</h2>
          </header>
          <List data={props.users} component={UserItem} />
        </div>
      </div>
      <div className="container">
        <div className="row">
          <header className="col-md-12">
            <hr className="yellow" />
            <h2 className="yellow">Primers</h2>
          </header>
          <List data={props.primers} component={PrimerItem} />
        </div>
      </div>
      <div className="container">
        <div className="row">
          <header className="col-md-12">
            <hr className="yellow" />
            <h2 className="yellow">Sources</h2>
          </header>
          <List data={props.sources} component={SourceItem} />
        </div>
      </div>
      <div className="container">
        <div className="row">
          <header className="col-md-12">
            <hr className="yellow" />
            <h2 className="yellow">Urls</h2>
          </header>
          <List data={props.urls} component={UrlItem} />
        </div>
      </div>
      <div className="container">
        <div className="row">
          <header className="col-md-12">
            <hr className="yellow" />
            <h2 className="yellow">Content</h2>
          </header>
          <List data={props.content} component={ContentItem} />
        </div>
      </div>
    </div>
  );
};

Stylesheet.propTypes = {
};

export default connect(() => {
  return {
    users: [
      { id: "4254923e-a3a3-4d6f-ae1d-f6af177215db", created: 1484086940, updated: 1484086940, username: "steve", type: "none", email: "steve@qri.io", name: "Steven", description: "", home_url: "" },
      { id: "9080d2c6-4e2a-4f5c-a49f-97dc1c6e1ecb", created: 1484080475, updated: 1484080475, username: "test_one", type: "none", email: "brendanobrienesq@gmail.com", name: "brendan", description: "", home_url: "" },
    ],
    primers: [
      {
        "id": "09364a0f-6c08-47a6-8618-29f85c30001c",
        "created": "2017-04-05T00:00:01Z",
        "updated": "2017-04-05T19:01:43Z",
        "shortTitle": "HUD",
        "title": "US Department of Housing & Urban Development",
        "description": "HUD's mission is to create strong, sustainable, inclusive communities and quality affordable homes for all. HUD is working to strengthen the housing market to bolster the economy and protect consumers; meet the need for quality affordable rental homes; utilize housing as a platform for improving quality of life; build inclusive and sustainable communities free from discrimination, and transform the way HUD does business.",
        "parent": null,
        "meta": null,
        "stats": {
          "urlCount": 46,
          "contentUrlCount": 0,
          "contentMetadataCount": 0
        }
      },
      {
        "id": "5b1031f4-38a8-40b3-be91-c324bf686a87",
        "created": "2017-01-01T00:00:01Z",
        "updated": "2017-04-05T19:01:43Z",
        "shortTitle": "",
        "title": "Environmental Protection Agency",
        "description": "The mission of the Environmental Protection Agency is to protect human health and the environment through the development and enforcement of regulations. The EPA is responsible for administering a number of laws that span various sectors, such as agriculture, transportation, utilities, construction, and oil and gas. In the budget for FY 2017, the agency lays out goals to better support communities and address climate change following the President’s Climate Action Plan. Additionally, the agency aims to improve community water infrastructure, chemical plant safety, and collaborative partnerships among federal, state, and tribal levels.",
        "parent": null,
        "meta": null,
        "stats": {
          "urlCount": 2314,
          "contentUrlCount": 1,
          "contentMetadataCount": 1
        }
      },
      {
        "id": "d9deff9d-15e8-43f1-9d00-51160c0bffbe",
        "created": "2017-01-01T00:00:01Z",
        "updated": "2017-04-05T19:01:43Z",
        "shortTitle": "",
        "title": "US Census",
        "description": "US Census Bureau",
        "parent": null,
        "meta": null,
        "stats": {
          "urlCount": 6320,
          "contentUrlCount": 28,
          "contentMetadataCount": 2
        }
      },
    ],
    sources: [
      {
        "id": "a9c3127f-3462-4534-a39f-179fb4d1fffd",
        "created": "2017-04-05T00:00:01Z",
        "updated": "2017-04-05T19:01:43Z",
        "title": "Fair Market Rents",
        "description": "DESCRIPTION",
        "url": "www.huduser.org/portal/datasets/fmr.html",
        "primerId": "09364a0f-6c08-47a6-8618-29f85c30001c",
        "crawl": true,
        "staleDuration": 43200000000000,
        "lastAlertSent": null,
        "meta": null,
        "stats": {
          "urlCount": 46,
          "contentUrlCount": 0,
          "contentMetadataCount": 0
        }
      },
      {
        "id": "590e001b-7060-4e54-bc81-c20c305a8155",
        "created": "2017-04-25T00:00:01Z",
        "updated": "2017-04-05T19:01:43Z",
        "title": "Hazardous Air Pollutants",
        "description": "Office of Air and Radiation",
        "url": "www.epa.gov/haps",
        "primerId": "5b1031f4-38a8-40b3-be91-c324bf686a87",
        "crawl": false,
        "staleDuration": 43200000000000,
        "lastAlertSent": null,
        "meta": null,
        "stats": {
          "urlCount": 51,
          "contentUrlCount": 0,
          "contentMetadataCount": 0
        }
      },
      {
        "id": "326fcfa0-d3e6-4b2d-8f95-e77220e16109",
        "created": "2017-01-01T00:00:01Z",
        "updated": "2017-04-07T19:34:13Z",
        "title": "epa.gov",
        "description": "entire epa site",
        "url": "www.epa.gov",
        "primerId": "5b1031f4-38a8-40b3-be91-c324bf686a87",
        "crawl": false,
        "staleDuration": 43200000000000,
        "lastAlertSent": null,
        "meta": null,
        "stats": {
          "urlCount": 13556,
          "contentUrlCount": 1,
          "contentMetadataCount": 1
        }
      },
    ],
    urls: [
      {
        "id": "3be66ee1-0265-40b2-baf8-611ebb5ff80c",
        "url": "https://www.epa.gov/caddis-vol2",
        "created": "2017-03-16T13:28:06Z",
        "updated": "2017-03-22T10:20:02Z",
        "lastGet": "2017-03-22T10:20:02.157868Z",
        "status": 401,
        "contentType": "text/html; charset=UTF-8",
        "contentSniff": "text/plain; charset=utf-8",
        "headers": [
          "Server",
          "Apache",
          "X-Content-Type-Options",
          "nosniff",
          "Www-Authenticate",
          "Basic realm=\"To review this site, please use the credentials supplied to you.\"",
          "Date",
          "Wed, 22 Mar 2017 10:20:03 GMT",
          "Strict-Transport-Security",
          "max-age=31536000; preload;",
          "X-Ua-Compatible",
          "IE=Edge,chrome=1",
          "Content-Length",
          "0",
          "Content-Type",
          "text/html; charset=UTF-8",
          "Connection",
          "keep-alive"
        ],
        "hash": "1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
      },
      {
        "id": "b5e261a5-cdaa-4e2d-8ad2-5928b05b31be",
        "url": "https://www.epa.gov/caddis-vol4",
        "created": "2017-03-16T13:28:06Z",
        "updated": "2017-03-22T10:20:07Z",
        "lastGet": "2017-03-22T10:20:07.261072Z",
        "status": 401,
        "contentType": "text/html; charset=UTF-8",
        "contentSniff": "text/plain; charset=utf-8",
        "headers": [
          "Www-Authenticate",
          "Basic realm=\"To review this site, please use the credentials supplied to you.\"",
          "Content-Length",
          "0",
          "Date",
          "Wed, 22 Mar 2017 10:20:08 GMT",
          "Connection",
          "keep-alive",
          "Strict-Transport-Security",
          "max-age=31536000; preload;",
          "Server",
          "Apache",
          "X-Ua-Compatible",
          "IE=Edge,chrome=1",
          "Content-Type",
          "text/html; charset=UTF-8",
          "X-Content-Type-Options",
          "nosniff"
        ],
        "hash": "1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
      },
      {
        "id": "b364ea23-0c88-4082-89bf-c6347d47486e",
        "url": "https://www.epa.gov/sites/production/files/2017-02/",
        "created": "2017-03-15T17:31:17Z",
        "updated": "2017-03-22T13:15:46Z",
        "title": "CADDIS"
      },
      {
        "id": "a26e6ba2-bee0-4f9d-b880-9a13e7531ffe",
        "url": "https://www.epa.gov/cameo",
        "created": "2017-03-16T13:28:06Z",
        "updated": "2017-03-16T13:28:06Z"
      },
    ],
    content: [
      { url: "", },
    ],
  };
})(Stylesheet);
