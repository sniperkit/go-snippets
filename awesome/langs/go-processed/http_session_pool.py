import random
from contextlib import contextmanager

from typing import List

import requests
from requests.sessions import Session
from requests.adapters import HTTPAdapter

from proxy_provider import ProxyProvider

USER_AGENT = ("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4)"
              " AppleWebKit/537.36 (KHTML, like Gecko)"
              " Chrome/58.0.3029.96 Safari/537.36")


class RealWorldAdapter(HTTPAdapter):

    def __init__(self, *args, **kwargs):
        self.timeout = kwargs['timeout']
        del kwargs['timeout']
        super(RealWorldAdapter, self).__init__(*args, **kwargs)

    def send(self, *args, **kwargs):
        kwargs['timeout'] = self.timeout
        return super(RealWorldAdapter, self).send(*args, **kwargs)


class RealWorldSessions:

    def __init__(self, num_sessions=80, agent=USER_AGENT, max_retries=20, timeout=10,
                 proxy_provider: ProxyProvider = None):
        self.num_sessions = num_sessions
        self.agent = agent
        self.max_retries = max_retries
        self.timeout = timeout
        self.proxy_provider = proxy_provider
        self.sessions = []  # type: List[Session]
        self.http_adapter = RealWorldAdapter(max_retries=max_retries, timeout=timeout)

        for i in range(num_sessions):
            session = self._create_session()
            self.sessions.append(session)

    def _create_session(self) -> Session:
        session = requests.session()
        if self.proxy_provider:
            session.proxies = self.proxy_provider.get_proxy()
        session.mount("http://", self.http_adapter)
        session.mount("https://", self.http_adapter)
        session.headers = {"user-agent": self.agent}

        return session

    @contextmanager
    def acquire(self) -> Session:
        i = random.randrange(0, self.num_sessions)
        try:

            yield self.sessions[i]
        except (requests.RequestException, ConnectionError):
            try:
                self.sessions[i].close()
            finally:
                self.sessions[i] = self._create_session()

        except IndexError:
            print("i = {}, sessions = {}, num_sessions = {}".format(i, len(self.sessions), self.num_sessions))
            quit()