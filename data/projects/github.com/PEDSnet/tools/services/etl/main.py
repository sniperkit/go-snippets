import os
import sys
import json
import uuid
import logging
import provenance
from datetime import datetime
from requests.exceptions import HTTPError
from flask import Flask, Response, request
from flask.ext.cors import CORS
from logger import logger
from gitutil import get_content, get_commits, get_commit, dedupe_commits
from cmputil import Changelog
from version import __version__

# Token for GitHub authorization. Set on startup.
GITHUB_AUTH_TOKEN = None

# URL template for a specific blob (ref + path)
GITHUB_BLOB_URL = 'https://github.com/PEDSnet/Data_Models/blob/{}/{}'


def json_defaults(o):
    if isinstance(o, provenance.entity):
        return o.json()

    raise TypeError


def ldjson(prov):
    for msg in prov:
        yield json.dumps(msg.json())
        yield '\n'


def generate_all_provenance(paths, model_name):
    prov = []

    for path in paths:
        for commit in get_commits(path, token=GITHUB_AUTH_TOKEN):
            try:
                model = get_content(commit['file_path'],
                                    token=GITHUB_AUTH_TOKEN,
                                    ref=commit['sha'])
                model['name'] = model_name
            except Exception:
                continue

            # Likely a 404 because the file moved in the current commit.
            if model is None:
                continue

            model['name'] = model_name

            prov.extend(provenance.generate(file_name=path,
                                            domain='pedsnet.etlconv',
                                            model=model,
                                            commit=commit))

    return prov


class Resource():
    def __init__(self, model_name, file_paths, versions=None):
        self.model_name = model_name
        self.file_paths = file_paths
        self.versions = versions or {}

    @property
    def current_file_path(self):
        "Returns the current path of the document in the repository."
        return self.file_paths[-1]

    def get_all_commits(self):
        "Gets all commits for all paths for this document."
        commits = []

        for path in self.file_paths:
            commits.extend(get_commits(path, token=GITHUB_AUTH_TOKEN))

        return dedupe_commits(commits)

    def _get_ref_path(self, ref):
        "Resolve the full ref and file path of the passed ref."
        for c in self.get_all_commits():
            if c['sha'].startswith(ref):
                return c['sha'], c['file_path']

        return 'master', self.current_file_path

    def _get_ref_path_at_time(self, ts):
        ref = None
        path = None

        for c in self.get_all_commits():
            if c['timestamp'] > ts:
                return ref, path

            ref = c['sha']
            path = c['file_path']

        return None, None

    def _get_version_ref_path(self, ver):
        "Gets the the ref path for a specific version."
        if ver in self.versions:
            return self._get_ref_path(self.versions[ver])

        return None, None

    def _request_ref_path(self):
        "Get the ref and path from request args."
        ref = request.args.get('ref')
        asof = request.args.get('asof')
        version = request.args.get('version')

        # Explicit ref provided.
        if ref:
            return self._get_ref_path(ref)

        # Version
        if version:
            return self._get_version_ref_path(version)

        # Time provided.
        if asof:
            ts = provenance.parse_date(asof)
            return self._get_ref_path_at_time(ts)

        # Default to the latest.
        return None, self.current_file_path

    def serve_document(self):
        "HTTP handler for the document."
        ref, path = self._request_ref_path()

        if not path:
            return 'Not found', 404

        try:
            commit = get_commit(path, token=GITHUB_AUTH_TOKEN, ref=ref)
        except HTTPError as e:
            return str(e), 503

        # Redirect to source on github.
        if request.args.get('r'):
            url = GITHUB_BLOB_URL.format(commit['sha'], path)
            return '', 302, {'Location': url}

        try:
            model = get_content(path, token=GITHUB_AUTH_TOKEN, ref=ref)
        except HTTPError as e:
            return str(e), 503

        if not model:
            return 'Not found', 404

        model['name'] = self.model_name

        content = json.dumps({
            'commit': {
                'sha': commit['sha'],
                'timestamp': commit['timestamp'],
                'date': commit['commit']['committer']['date'],
                'file_path': commit['file_path'],
            },
            'model': model,
        })

        resp = Response(content)

        resp.headers['Content-Type'] = 'application/json'
        return resp

    def serve_commits(self):
        commits = []

        for c in self.get_all_commits():
            commits.append({
                'file_path': c['file_path'],
                'sha': c['sha'],
                'date': c['commit']['committer']['date'],
            })

        resp = Response(json.dumps(commits))
        resp.headers['content-type'] = 'application/json'

        return resp

    def serve_provenance(self):
        "HTTP handler for the provenance."
        ref, path = self._request_ref_path()

        if not path:
            return 'not found', 404

        try:
            model = get_content(path, token=GITHUB_AUTH_TOKEN, ref=ref)
        except HTTPError as e:
            return str(e), 503

        if not model:
            return 'not found', 404

        model['name'] = self.model_name

        try:
            commit = get_commit(path, token=GITHUB_AUTH_TOKEN, ref=ref)
        except HTTPError as e:
            return str(e), 503

        prov = provenance.generate(file_name=self.current_file_path,
                                   domain='pedsnet.etlconv',
                                   model=model,
                                   commit=commit)

        if request.accept_mimetypes.best == 'application/json; boundary=NL':
            content = ldjson(prov)
            content_type = 'application/json; boundary=NL'
        else:
            content = json.dumps(prov, default=json_defaults)
            content_type = 'application/json'

        resp = Response(content)
        resp.headers['content-type'] = content_type

        return resp

    def serve_full_provenance(self):
        prov = generate_all_provenance(self.file_paths, self.model_name)

        if request.accept_mimetypes.best == 'application/json; boundary=NL':
            content = ldjson(prov)
            content_type = 'application/json; boundary=NL'
        else:
            content = json.dumps(prov, default=json_defaults)
            content_type = 'application/json'

        resp = Response(content)
        resp.headers['content-type'] = content_type

        return resp

    def serve_changes_all(self):
        "HTTP handler for serving the change log."
        prov = generate_all_provenance(self.file_paths, self.model_name)

        cl = Changelog()

        log = []

        for e in prov:
            if not {'Model', 'Table', 'Field'} & set(e.labels):
                continue

            c = cl.evaluate(e.json())

            if c is not None:
                log.append(c)

        if request.accept_mimetypes.best == 'application/json; boundary=NL':
            content = ldjson(log)
            content_type = 'application/json; boundary=NL'
        else:
            content = json.dumps(log, default=json_defaults)
            content_type = 'application/json'

        resp = Response(content)
        resp.headers['content-type'] = content_type

        return resp


pedsnet = Resource(
        model_name='pedsnet',
        file_paths=(
            'PEDSnet/docs/PEDSnet_CDM_V1_ETL_Conventions.md',
            'PEDSnet/V1/docs/PEDSnet_CDM_V1_ETL_Conventions.md',
            'PEDSnet/V2/docs/Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md',
            'PEDSnet/docs/Pedsnet_CDM_ETL_Conventions.md',
        ), versions={
            '1.0.0': 'ad18c4ea1e227bbccf3fbc0a5ae05b1f552af95d',
            '2.0.0': '530d08afdff1542fcbc9042794a90e9e444541c7',
            '2.1.0': 'ed80cfd20fca8e2c3495cbd0b2971c1d713f6be6',
            '2.2.0': 'master',
        })

i2b2 = Resource(
        model_name='i2b2',
        file_paths=(
            'i2b2/V2/docs/i2b2_pedsnet_v2_etl_conventions.md',
        ), versions={
            '2.0.0': 'master',
        })


# Initialize the flask app and register the routes.
app = Flask(__name__)
CORS(app)

# Unique service ID and timestamp when it started.
SERVICE_ID = str(uuid.uuid4())
SERVICE_TIMESTAMP = datetime.now().strftime('%Y-%m-%dT%H:%M:%S')


@app.route('/', methods=['GET'])
def index():
    return json.dumps({
        'name': 'ETL Conventions Service',
        'version': __version__,
        'time': SERVICE_TIMESTAMP,
        'uuid': SERVICE_ID,
    })

app.add_url_rule('/pedsnet',
                 'pedsnet_document',
                 pedsnet.serve_document,
                 methods=['GET'])

app.add_url_rule('/pedsnet/commits',
                 'pedsnet_commits',
                 pedsnet.serve_commits,
                 methods=['GET'])

app.add_url_rule('/pedsnet/prov',
                 'pedsnet_prov',
                 pedsnet.serve_provenance,
                 methods=['GET'])

app.add_url_rule('/pedsnet/prov/all',
                 'pedsnet_prov_all',
                 pedsnet.serve_full_provenance,
                 methods=['GET'])

app.add_url_rule('/pedsnet/prov/changes/all',
                 'pedsnet_changes_all',
                 pedsnet.serve_changes_all,
                 methods=['GET'])

app.add_url_rule('/i2b2',
                 'i2b2_document',
                 i2b2.serve_document,
                 methods=['GET'])

app.add_url_rule('/i2b2/commits',
                 'i2b2_commits',
                 i2b2.serve_commits,
                 methods=['GET'])

app.add_url_rule('/i2b2/prov',
                 'i2b2_prov',
                 i2b2.serve_provenance,
                 methods=['GET'])

app.add_url_rule('/i2b2/prov/all',
                 'i2b2_prov_all',
                 i2b2.serve_full_provenance,
                 methods=['GET'])

app.add_url_rule('/i2b2/prov/changes/all',
                 'i2b2_changes_all',
                 i2b2.serve_changes_all,
                 methods=['GET'])


if __name__ == '__main__':
    usage = """PEDSnet ETL Conventions Service

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

    app.run(host=host, port=port, threaded=True, debug=debug)
