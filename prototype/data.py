"""This file provides abstraction over the source of data used in the server"""

import random
import json
import datetime
import codecs


class DataSource:
    """This class abstracts the underlying system that stores the computed
    features about the developers"""
    def __init__(self):
        pass


class RandomData(DataSource):
    def __init__(self, seed=None):
        r = random.Random(seed)
        # print(r.random())

        self.scores = {}
        self.features = ['python', 'cookies', 'chocolate',
                         'superpowers', 'doge']
        self.developers = ['Robin', 'Kevin', 'Laurent', 'Xuan',
                           'Frederik', 'Daniel', 'Clément']

        for developer in self.developers:
            scores = {}
            for feature in self.features:
                scores[feature] = r.random()
            self.scores[developer] = scores


def iter_repos(repo_file):
    with codecs.open(repo_file, 'r', 'utf-8') as f:
        for repo in f:
            yield json.loads(repo)
        return


class JsonData(DataSource):
    def __init__(self, dev_file, repo_file):

        devs = []
        with codecs.open(dev_file, 'r', 'utf-8') as f:
            for dev in f:
                devs.append(json.loads(dev))

        self.developers = []
        self.scores = {}

        for dev in devs:
            self.developers.append(dev['login'])
            self.scores[dev['login']] = self.compute_scores(dev)

        # scores related to repos
        self.compute_scores_repos(repo_file)

        # Compute the list of features
        self.features = set()
        for v in self.scores.values():
            for k in v.keys():
                self.features.add(k.replace(' ', ''))
        self.features = list(self.features)

    def compute_scores_repos(self, repo_file):
        for repo in iter_repos(repo_file):
            owner = repo['owner']['login']
            if owner in self.scores.keys():
                score = self.scores[owner]
                if repo['language'] is not None:
                    if repo['language'].lower() in score.keys():
                        score[repo['language'].lower()] += int(repo['size'])
                    else:
                        score['number_languages'] += 1
                        score[repo['language'].lower()] = int(repo['size'])

    def compute_scores(self, dev):
        score = {}

        score['followers'] = dev['followers']
        score['public_repos'] = dev['public_repos']

        created = datetime.datetime.strptime(dev['created_at'],
                                             "%Y-%m-%dT%H:%M:%SZ")
        updated = datetime.datetime.strptime(dev['updated_at'],
                                             "%Y-%m-%dT%H:%M:%SZ")
        score['created_at'] = int(created.timestamp())
        score['updated_at'] = int(updated.timestamp())

        score['number_languages'] = 0

        return score
