import io
import os
import requests
from threading import Lock
from parser import Parser
from logger import logger
from dateutil import parse_date

# Required mediatype for the Accept header to get the raw content.
GITHUB_RAW_MEDIATYPE = 'application/vnd.github.v3.raw'

# Default mediatype for GitHub responses.
GITHUB_REQUEST_MEDIATYPE = 'application/vnd.github.v3'

# Timeout for the request.
REQUEST_TIMEOUT = 20

# The commits endpoint is used to get the commits for a particular file.
GITHUB_COMMITS_URL = 'https://api.github.com/repos/PEDSnet/Data_Models/commits'

# The contents endpoint is used to get the contents of a file at
# a paritcular revision.
GITHUB_CONTENTS_URL = 'https://api.github.com/repos/PEDSnet/Data_Models/contents'  # noqa


def dedupe_commits(commits):
    shas = set()
    filtered = []

    # Evaluate in descending order so the more recent commit is included
    # in the set.
    for c in reversed(commits):
        if c['sha'] in shas:
            continue

        filtered.append(c)
        shas.add(c['sha'])

    filtered.reverse()
    return filtered


# Get all commits for the file.
def get_commits(path, token):
    headers = {
        'Accept': GITHUB_REQUEST_MEDIATYPE,
        'Authorization': 'token ' + token,
    }

    resp = requests.get(GITHUB_COMMITS_URL,
                        params={'path': path},
                        headers=headers,
                        timeout=REQUEST_TIMEOUT)

    resp.raise_for_status()

    # Commits are in descending order, so we reverse them.
    commits = resp.json()
    commits.reverse()

    # Add the file path to the annotation
    for i, c in enumerate(commits):
        commits[i] = {
            'sha': c['sha'],
            'commit': c['commit'],
            'file_path': path,
            'timestamp': parse_date(c['commit']['committer']['date']),  # noqa
        }

    return commits


def get_commit(path, token, ref='master'):
    headers = {
        'Accept': GITHUB_REQUEST_MEDIATYPE,
        'Authorization': 'token ' + token,
    }

    resp = requests.get(GITHUB_COMMITS_URL,
                        params={
                            'path': path,
                            'sha': ref,
                        },
                        headers=headers,
                        timeout=REQUEST_TIMEOUT)

    resp.raise_for_status()
    c = resp.json()[0]

    return {
        'sha': c['sha'],
        'commit': c['commit'],
        'file_path': path,
        'timestamp': parse_date(c['commit']['committer']['date']),
    }


# Cache for file content. The key is a tuple of the path and ref.
class ContentCache():
    def __init__(self):
        self.cache = {}
        self.locks = {}

        # Lock to define new locks.
        self._lock = Lock()

    def __contains__(self, key):
        return key in self.cache

    def __getitem__(self, key):
        return self.cache.get(key)

    def __setitem__(self, key, value):
        self.cache[key] = value

    def lock(self, key):
        self._lock.acquire()

        if key in self.locks:
            lock = self.locks[key]
        else:
            lock = Lock()
            self.locks[key] = lock

        lock.acquire(timeout=3)
        self._lock.release()

    def unlock(self, key):
        self.locks[key].release()


content_cache = ContentCache()


def get_content(path, token, ref=None):
    key = (path, ref)

    if ref and key in content_cache:
        text = content_cache[key]
        logger.debug('[content] cache hit %s', key)
    else:
        logger.debug('[content] cache miss %s', key)

        if ref:
            # Lock to prevent other
            content_cache.lock(key)

        headers = {
            'Accept': GITHUB_RAW_MEDIATYPE,
            'Authorization': 'token ' + token,
        }

        document_url = os.path.join(GITHUB_CONTENTS_URL, path)

        resp = requests.get(document_url,
                            params={'ref': ref},
                            headers=headers,
                            timeout=REQUEST_TIMEOUT)

        if resp.status_code == 404:
            text = None
        else:
            resp.raise_for_status()
            text = resp.text

        if ref:
            # Update cache and release the lock.
            content_cache[key] = text
            content_cache.unlock(key)

    if text is None:
        return

    # Wrap decoded bytes in file-like object.
    buff = io.StringIO(text)

    return Parser(buff).parse()
