import threading


class ResourceTimer():
    def __init__(self, resource, interval):
        "interval is in minutes"
        if interval < 1:
            interval = 10
        interval *= 60

        self.interval = interval
        self.resource = resource

        self.thread = threading.Timer(self.interval, self.update)

    def update(self):
        self.resource.update()
        self.thread = threading.Timer(self.interval, self.update)
        self.thread.start()

    def start(self):
        self.thread.start()

    def cancel(self):
        self.thread.cancel()
