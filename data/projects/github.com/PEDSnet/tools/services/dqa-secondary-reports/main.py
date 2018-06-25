import codecs
from collections import defaultdict
from datetime import datetime
import git
import json
import logging
import os
import shutil
import sys
import uuid

from flask import Flask, Response
from flask.ext.cors import CORS

from dictionary_parser import Parser as DictionaryParser
from report_parser import Parser, logger
from timer import ResourceTimer
from version import __version__

# Token for GitHub authorization. Set on startup.
GITHUB_AUTH_TOKEN = None

REPO_DIR = 'dqa_repo'
REMOTE_URL = 'github.com/PEDSnet/Data-Quality.git'

# Recursive nested dict structure.
rdict = lambda: defaultdict(rdict)


class DQAResource():
    def __init__(self, dir_name):
        self.field_totals = rdict()
        self.table_totals = rdict()
        self.site_table_totals = rdict()
        self.site_totals = rdict()
        self.expanded_site_totals = rdict()
        self.site_list = rdict()
        self.dictionary = {}

        self.dir_name = dir_name
        self.remote_url = 'https://' + GITHUB_AUTH_TOKEN + '@' + REMOTE_URL

    def process_dqa(self):
        self.field_totals = rdict()
        self.table_totals = rdict()
        self.site_table_totals = rdict()
        self.site_totals = rdict()
        self.expanded_site_totals = rdict()
        self.site_list = rdict()

        dir_name = self.dir_name + '/SecondaryReports/'

        for site in os.listdir(dir_name):
            hosp_path = os.path.join(dir_name, site)

            if site == 'Ranking' or not os.path.isdir(hosp_path):
                continue

            # for each site, we only need data from the latest report for each
            # data model version, so go through the dqa reports in reverse
            # order and ignore older reports for each data model version.
            for etl_version in sorted(os.listdir(hosp_path), reverse=True):
                etl_path = os.path.join(hosp_path, etl_version)

                if etl_version.startswith('ETL') and os.path.isdir(etl_path):
                    file_names = []
                    for name in os.listdir(etl_path):
                        if name.endswith('.csv'):
                            file_names.append(name)

                    # there shouldn't be any special characters in the
                    # DQA report files, but some of the older reports do
                    # contain latin-1 characters.
                    version = Parser.get_version(
                        codecs.open(os.path.join(etl_path, file_names[0]), 'r',
                                    encoding='latin1'))

                    if site in self.site_totals[version]:
                        continue

                    for file_name in file_names:
                        table = file_name.split('.csv')[0]
                        buff = codecs.open(
                            os.path.join(etl_path, file_name), 'r',
                            encoding='latin1')
                        p = Parser(buff, site, table,
                                   self.field_totals[version],
                                   self.table_totals[version],
                                   self.site_table_totals[version][site],
                                   self.site_totals[version][site],
                                   self.expanded_site_totals[version][site])
                        p.parse()

                        if table in self.site_table_totals[version][site]:
                            if table not in self.site_list[version]:
                                self.site_list[version][table] = [site]
                            else:
                                self.site_list[version][table].append(site)

        for version in self.site_totals:
            self.site_list[version]['all'] = []

            for site in self.site_totals[version]:
                self.site_list[version]['all'].append(site)

        # if a site has no issues with a certain status,
        # add a corresponding entry and set it to 0.
        for version in self.site_totals:
            status_names = []
            for site in self.site_totals[version]:
                for status in self.site_totals[version][site]:
                    status_names.append(status)

            status_names = set(status_names)

            for site in self.site_totals[version]:
                for status in status_names:
                    if status not in self.site_totals[version][site]:
                        self.site_totals[version][site][status] = 0

                    for table in self.expanded_site_totals[version][site]:
                        if status not in self.expanded_site_totals[version][site][table]:
                            self.expanded_site_totals[version][site][table][status] = 0

    def process_dictionary(self):
        path = self.dir_name + '/Dictionary/DCC_DQA_Dictionary.csv'
        buff = codecs.open(path, 'r', encoding='latin1')
        self.dictionary = DictionaryParser(buff).parse()

    def update(self):
        self.update_repo()
        self.process_dqa()
        self.process_dictionary()

    def update_repo(self):
        logger.debug('start updating repo')

        needs_refresh = True

        if os.path.isdir(self.dir_name):
            try:
                git.Repo(self.dir_name).remotes.origin.pull()
                needs_refresh = False
            except:
                shutil.rmtree(self.dir_name)

        if needs_refresh:
            os.mkdir(self.dir_name)
            git.Repo.clone_from(self.remote_url, self.dir_name)

        logger.debug('done updating repo')

    def serve_field_totals(self, version, table=None, field=None):
        "Entrypoint for Flask routing for a given version of PEDSnet data model"
        def versioned_data():
            if table is None:
                info = self.field_totals[version]
            elif field is None:
                info = self.field_totals[version][table]
            else:
                info = self.field_totals[version][table][field]

            return wrap_response(info)

        return versioned_data

    def serve_dict(self):
        "Entrypoint for Flask routing for dictionary of DQA status codes"

        return wrap_response(self.dictionary)

    def serve_site_totals(self, version, site=None):
        """
        Entrypoint for Flask routing for summary DQA report that lists,
        for every site, the number of issues of each status (e.g. new,
        persistent, etc).
        The reposrt is given for a given version of PEDSnet data model.
        """
        def versioned_data():
            if site is None:
                info = self.site_totals[version]
            else:
                info = self.site_totals[version][site]

            return wrap_response(info)

        return versioned_data

    def serve_expanded_site_totals(self, version, site):
        """
        Entrypoint for Flask routing for summary DQA report for a specific site
        and a given version of PEDSnet data model.
        The report lists the number of issues of each status raised for this site
        for each of the tatbles in the model.
        """
        def versioned_data():
            return wrap_response(self.expanded_site_totals[version][site])

        return versioned_data

    def serve_table_totals(self, version, table=None):
        "Entrypoint for Flask routing for table totals for a given version of PEDSnet data model"
        def versioned_data():
            if table is None:
                info = self.table_totals[version]
            else:
                info = self.table_totals[version][table]

            return wrap_response(info)

        return versioned_data

    def serve_site_list(self, version, table=None):
        """
        Entrypoint for Flask routing that provides a list of sites that have
        an outstanding issue against the given table in a given version of the
        PEDSnet data model.
        """
        def versioned_data():
            if table is None:
                info = self.site_list[version]['all']
            else:
                info = self.site_list[version][table]

            return wrap_response(info)

        return versioned_data

    def serve_site_table_totals(self, version, table, site):
        """
        Entrypoint for Flask routing for (site+table)-specific totals
        for a given version of PEDSnet data model
        """
        def versioned_data():
            info = self.site_table_totals[version][site][table]

            return wrap_response(info)

        return versioned_data

    def get_versions(self):
        return list(self.field_totals.keys())

    def serve_versions(self):
        "Entrypoint for Flask routing for PEDSnet data model versions"

        return wrap_response(self.get_versions())


def wrap_response(info):
    content = json.dumps(info, indent=4, sort_keys=True)

    resp = Response(content)

    resp.headers['Content-Type'] = 'application/json'

    return resp

if __name__ == '__main__':
    usage = """PEDSnet DQA Secondary Reports Service

    Usage: main.py [--debug]
                   [--host=<host>] [--port=<port>]
                   [--token=<token>]
                   [--interval=<interval>]
                   [--repo-dir=<repo-dir>]

    Options:
        --help                  Display the help.
        --debug                 Enable debug output.
        --host=<host>           Host of the service.
        --port=<port>           Port of the service [default: 5000].
        --token=<token>         GitHub authorization token.
        --interval=<interval>   Time interval, in minutes, to periodically check the repo for updates [default: 30]
        --repo-dir=<repo-dir>   Directory of the repo on disk.
    """  # noqa

    from docopt import docopt

    opts = docopt(usage)

    host = opts['--host']
    port = int(opts['--port'])
    debug = opts['--debug']
    interval = int(opts['--interval'])

    # Set token if defined.
    if opts['--token']:
        GITHUB_AUTH_TOKEN = opts['--token']
    else:
        GITHUB_AUTH_TOKEN = os.environ.get('GITHUB_AUTH_TOKEN')

    if not GITHUB_AUTH_TOKEN:
        print('Authorization token required.')
        sys.exit(1)

    if opts['--repo-dir']:
        repo_dir = opts['--repo-dir']
    elif 'REPO_DIR' in os.environ:
        repo_dir = os.environ['REPO_DIR']
    else:
        repo_dir = REPO_DIR

    if debug:
        logger.setLevel(logging.DEBUG)

    dqa_resource = DQAResource(repo_dir)

    # Do the initial fetch of the repository.
    dqa_resource.update()

    # Create the timer for repo update, and start it.
    timer = ResourceTimer(dqa_resource, interval)
    timer.start()

    # Initialize the flask app and register the routes.
    app = Flask(__name__)

    CORS(app)

    # Unique service ID and timestamp when it started.
    SERVICE_ID = str(uuid.uuid4())
    SERVICE_TIMESTAMP = datetime.now().strftime('%Y-%m-%dT%H:%M:%S')

    @app.route('/', methods=['GET'])
    def index():
        return json.dumps({
            'name': 'DQA Service',
            'version': __version__,
            'time': SERVICE_TIMESTAMP,
            'uuid': SERVICE_ID,
        })

    """
    Make sure all urls end with a slash. The version without a trailing slash
    will be taken care of automatically because Flask will automatically
    redirect the user to the same URL with a trailing slash attached.
    The reverse is not true.
    """

    app.add_url_rule('/pedsnet/',
                     'pedsnet_versions',
                     dqa_resource.serve_versions,
                     methods=['GET'])

    for version in dqa_resource.get_versions():
        urlf = '/pedsnet/{version}/field-totals/'
        namef = '{version}_field-totals'

        app.add_url_rule(urlf.format(version=version),
                         namef.format(version=version),
                         dqa_resource.serve_field_totals(version),
                         methods=['GET'])

        urlf = '/pedsnet/{version}/field-totals/{table}/'
        namef = '{version}_field-totals_{table}'

        for table in dqa_resource.field_totals[version]:
            app.add_url_rule(urlf.format(version=version, table=table),
                             namef.format(version=version, table=table),
                             dqa_resource.serve_field_totals(version, table),
                             methods=['GET'])

        urlf = '/pedsnet/{version}/field-totals/{table}/{field}/'
        namef = '{version}_field-totals_{table}_{field}'

        for table in dqa_resource.field_totals[version]:
            for field in dqa_resource.field_totals[version][table]:
                app.add_url_rule(urlf.format(version=version, table=table, field=field),
                                 namef.format(version=version, table=table, field=field),
                                 dqa_resource.serve_field_totals(version, table, field),
                                 methods=['GET'])

        urlf = '/pedsnet/{version}/site-totals/'
        namef = '{version}_site-totals'

        app.add_url_rule(urlf.format(version=version),
                         namef.format(version=version),
                         dqa_resource.serve_site_totals(version),
                         methods=['GET'])

        urlf = '/pedsnet/{version}/site-list/'
        namef = '{version}_site_list'

        app.add_url_rule(urlf.format(version=version),
                         namef.format(version=version),
                         dqa_resource.serve_site_list(version),
                         methods=['GET'])

        urlf = '/pedsnet/{version}/site-totals/{site}/'
        namef = '{version}_site-totals_{site}'

        for site in dqa_resource.site_totals[version]:
            app.add_url_rule(urlf.format(version=version, site=site),
                             namef.format(version=version, site=site),
                             dqa_resource.serve_expanded_site_totals(version, site),
                             methods=['GET'])

        urlf = '/pedsnet/{version}/table-totals/'
        namef = '{version}_table-totals'

        app.add_url_rule(urlf.format(version=version),
                         namef.format(version=version),
                         dqa_resource.serve_table_totals(version),
                         methods=['GET'])

        urlf = '/pedsnet/{version}/table-totals/{table}/'
        namef = '{version}_table-totals_{table}'

        for table in dqa_resource.table_totals[version]:
            app.add_url_rule(urlf.format(version=version, table=table),
                             namef.format(version=version, table=table),
                             dqa_resource.serve_table_totals(version, table),
                             methods=['GET'])

        urlf = '/pedsnet/{version}/site-list/{table}/'
        namef = '{version}_site_list_{table}'

        for table in dqa_resource.table_totals[version]:
            app.add_url_rule(urlf.format(version=version, table=table),
                             namef.format(version=version, table=table),
                             dqa_resource.serve_site_list(version, table),
                             methods=['GET'])

        urlf = '/pedsnet/{version}/table-totals/{table}/{site}/'
        namef = '{version}_table-totals_{table}_{site}'

        for table in dqa_resource.table_totals[version]:
            for site in dqa_resource.site_totals[version]:
                app.add_url_rule(urlf.format(version=version, table=table, site=site),
                                 namef.format(version=version, table=table, site=site),
                                 dqa_resource.serve_site_table_totals(version, table, site),
                                 methods=['GET'])

    app.add_url_rule('/dictionary/', 'dictionary', dqa_resource.serve_dict,
                     methods=['GET'])

    try:
        '''
        Set use_reloader=False when starting the app, otherwise we end up with two resourse update timers
        starting and then stepping on each other's toes.

        For explanation, see:
        http://stackoverflow.com/questions/25504149/why-does-running-the-flask-dev-server-run-itself-twice

        The Werkzeug reloader spawns a child process so that it can restart that process each time
        your code changes.
        Werkzeug is the library that supplies Flask with the development server when you call app.run().
        See the restart_with_reloader() function code; your script is run again with subprocess.call().

        If you set use_reloader to False you'll see the behaviour go away, but then you also lose the
        reloading functionality.
        '''
        app.run(host=host, port=port, debug=debug, use_reloader=False)
    finally:
        timer.cancel()
