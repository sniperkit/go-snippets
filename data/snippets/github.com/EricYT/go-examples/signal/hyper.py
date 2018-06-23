#!/usr/bin/env python

import time
import os
from subprocess import CalledProcessError, Popen


class Parent:
    def __init__(self):
        print("Parent init")

    def start(self):
        print("Ready to start command daemon")
        try:
            #Popen(['./signal'])
            os.system("./signal &")
        except CalledProcessError as err:
            print("Run command daemon error: {}".format(err.output))
            return

    def wait(self):
        print("Waiting to die")
        time.sleep(3)
        print("Main process died")

if __name__ == '__main__':
    parent = Parent()
    parent.start()
    parent.wait()
