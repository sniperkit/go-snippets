import codecs
import json
import time
from getpass import getpass

from github3 import (
    GitHub,
    GitHubError,
    login
)


import gh_settings as settings


def main():
    user = input('GitHub username: ')
    password = getpass('GitHub password: ')

    if user and password:
        gh = login(user, password)
    else:
        gh = GitHub()

    # It is defined inside the main because it makes use of the current gh
    def call_or_wait(function, *args, **kargs):
        '''Given a function and some arguments, tries to execute the function
        or waits until the ratelimit is reset.'''

        success = False
        while not success:
            try:
                result = function(*args, **kargs)
                success = True
            except GitHubError as e:
                if e.code == 403:
                    # rate limit exceeded
                    reset = gh.rate_limit()['rate']['reset']
                    wait_time = int(reset - time.time()) + 1

                    print("Not enough API calls. Waiting for",
                          int(wait_time / 60), "minutes and",
                          wait_time % 60, "seconds")

                    time.sleep(wait_time)
                else:
                    raise e

        return result

    # get awesome developers
    devmine = call_or_wait(lambda:
                           gh.organization("DevMine").iter_public_members())
    # the lambda is used to make the call lazy, as there are multiple API
    # calls involved in the statement

    devs = {dev.login for dev in devmine}

    # get random developers
    for repo in call_or_wait(gh.iter_all_repos,
                             number=settings.REPOS_COUNT,
                             since=settings.REPOS_SINCE_ID):
        # skip organizations and forked repository to find users that actually
        # added new content
        if repo.private or repo.fork or (repo.owner.type == 'Organization'):
            continue
        devs.add(repo.owner.login)

    print(len(devs), "developers fetched")

    fpu = codecs.open(settings.USERS_DATASET, 'a', 'utf-8')
    fpr = codecs.open(settings.REPOS_DATASET, 'a', 'utf-8')

    # dump developers and their repositories
    users_processed = 0
    for dev in devs:
        u = call_or_wait(gh.user, dev)
        json.dump(u.to_json(), fpu)
        fpu.write("\n")
        users_processed += 1
        if users_processed % 100 == 0:
            print("Repos feched for", users_processed, "developers")
        for repo in call_or_wait(gh.iter_user_repos, u.login):
            json.dump(repo.to_json(), fpr)
            fpr.write("\n")

    fpu.close()
    fpr.close()


if __name__ == "__main__":
    main()
