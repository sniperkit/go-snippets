import daemon
import time
from classes import SnapSet

if __name__ == "__main__":
	with daemon.DaemonContext():
		s = SnapSet()
		s.start()
