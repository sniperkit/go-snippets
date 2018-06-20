from requests.adapters import HTTPAdapter


class RealWorldAdapter(HTTPAdapter):

    def __init__(self, *args, **kwargs):
        self.timeout = kwargs['timeout']
        del kwargs['timeout']
        super(RealWorldAdapter, self).__init__(*args, **kwargs)

    def send(self, *args, **kwargs):
        kwargs['timeout'] = self.timeout
        return super(RealWorldAdapter, self).send(*args, **kwargs)
