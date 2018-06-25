import io
import os
import sys
import json
import uuid
import logging
import requests
from datetime import datetime
from requests.exceptions import HTTPError
from flask import Flask, Response
from flask.ext.cors import CORS
from parser import Parser, logger
from version import __version__


# Required mediatype for the Accept header to get the raw content.
GITHUB_RAW_MEDIATYPE = 'application/vnd.github.v3.raw'
GITHUB_DEFAULT_MEDIATYPE = 'application/vnd.github.v3'

# Token for GitHub authorization. Set on startup.
GITHUB_AUTH_TOKEN = None

DEFAULT_TIMEOUT = 5


class DataDictionaryResource():
    def __init__(self, document_url, commits_url, file_path):
        self.document_url = document_url
        self.commits_url = commits_url
        self.file_path = file_path

        self.content_last_modified = None
        self.content_etag = None

        self.commit_last_modified = None
        self.commit_etag = None

        self.cached_model = None
        self.cached_commit = None

    def __call__(self):
        "Entrypoint for Flask routing."
        try:
            model = self.parse_model()
        except HTTPError:
            return '', 503

        try:
            commit = self.parse_commit()
        except HTTPError:
            return '', 503

        content = json.dumps({
            'commit': commit,
            'model': model,
        })

        resp = Response(content)

        resp.headers['Content-Type'] = 'application/json'

        return resp

    def parse_model(self):
        headers = {
            'Accept': GITHUB_RAW_MEDIATYPE,
        }
        if GITHUB_AUTH_TOKEN:
            headers['Authorization'] = 'token ' + GITHUB_AUTH_TOKEN

        if self.content_last_modified:
            headers['If-Modified-Since'] = self.content_last_modified

        if self.content_etag:
            headers['If-None-Match'] = self.content_etag

        resp = requests.get(self.document_url,
                            headers=headers,
                            timeout=DEFAULT_TIMEOUT)

        resp.raise_for_status()

        self.content_last_modified = resp.headers['Last-Modified']
        self.content_etag = resp.headers['ETag']

        # Not modified based on the conditional headers.
        if resp.status_code == 200:
            # Wrap decoded bytes in file-like object.
            buff = io.BytesIO(resp.content)
            self.cached_model = Parser(buff).parse()

        return self.cached_model

    def parse_commit(self):
        headers = {
            'Accept': GITHUB_DEFAULT_MEDIATYPE,
        }
        if GITHUB_AUTH_TOKEN:
            headers['Authorization'] = 'token ' + GITHUB_AUTH_TOKEN

        if self.commit_last_modified:
            headers['If-Modified-Since'] = self.commit_last_modified

        if self.commit_etag:
            headers['If-None-Match'] = self.commit_etag

        resp = requests.get(self.commits_url,
                            params={'path': self.file_path},
                            headers=headers,
                            timeout=DEFAULT_TIMEOUT)

        resp.raise_for_status()

        self.commit_last_modified = resp.headers['Last-Modified']
        self.commit_etag = resp.headers['ETag']

        # Not modified based on the conditional headers.
        if resp.status_code == 200:
            # Get the most recent commit.
            commit = resp.json()[0]

            self.cached_commit = {
                'sha': commit['sha'],
                'date': commit['commit']['committer']['date']
            }

        return self.cached_commit


pcornet_v3 = DataDictionaryResource(
        document_url='https://api.github.com/repos/PEDSnet/Data_Coordinating_Center/contents/output/pcornet/annotated_data_dictionary.xlsx',  # noqa
        commits_url='https://api.github.com/repos/PEDSnet/Data_Coordinating_Center/commits',  # noqa
        file_path='output/pcornet/annotated_data_dictionary.xlsx')

# Initialize the flask app and register the routes.
app = Flask(__name__)
CORS(app)


# Unique service ID and timestamp when it started.
SERVICE_ID = str(uuid.uuid4())
SERVICE_TIMESTAMP = datetime.now().strftime('%Y-%m-%dT%H:%M:%S')


@app.route('/', methods=['GET'])
def index():
    return json.dumps({
        'name': 'PCORnet Annotated Data Dictionary Service',
        'version': __version__,
        'time': SERVICE_TIMESTAMP,
        'uuid': SERVICE_ID,
    })

# The URL should have a trailing slash. Flask will redirect a
# request without a trailing slash to the same URL with
# a trailing slash attached -- but the reverse is not true.
app.add_url_rule('/pcornet/3.0.0/',
                 'pcornet_v3',
                 pcornet_v3,
                 methods=['GET'])

if __name__ == '__main__':
    usage = """PCORnet Annotated Data Dictionary Service

    Usage: main.py [--token=<token>] [--host=<host>] [--port=<port>] [--debug]

    Options:
        --help              Display the help.
        --token=<token>     GitHub authorization token.
        --host=<host>       Host of the service.
        --port=<port>       Port of the service [default: 5000].
        --debug             Enable debug output.
    """

    from docopt import docopt

    opts = docopt(usage)

    host = opts['--host']
    port = int(opts['--port'])
    debug = opts['--debug']

    # Set token if defined.
    if opts['--token']:
        GITHUB_AUTH_TOKEN = opts['--token']
    else:
        GITHUB_AUTH_TOKEN = os.environ.get('GITHUB_AUTH_TOKEN')

    if not GITHUB_AUTH_TOKEN:
        print('Authorization token required.')
        sys.exit(1)

    if debug:
        logger.setLevel(logging.DEBUG)

    app.run(host=host, port=port, debug=debug)
