import daemon
import datetime
import logging
import os
import shutil
import tarfile
import time
from config import SnapConfig
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler

class SnapSetHandler(FileSystemEventHandler):
	def __init__(self):
		self.s = SnapSet()
		self.s.start()
	def on_any_event(self, event):
		print event.src_path
		#self.s.add_source(event.src_path)
		if self.s.first == False and self.s.done == True:
			self.s.start()

class SnapSet(object):
	'''
	Back up changes in real time to multipal destinations
	'''
	dict_free_size = {}
	def __init__(self, src=None, dest=None, interval=None, times=None,
		     git=False, compress=False, dryrun=False):
		self.first = True
		self.done = False
		self.timestamp = self.get_timestamp()
		if src != None and dest != None and interval != None and times != None:
			self.src = src
			self.dest = dest
			self.interval = interval
			self.times = times
			self.git = git
			self.compress = compress
		else:
			self.src, self.dest, self.interval, self.times, self.git, self.compress = self._load_config("./test.ini")

	def _load_config(self, name): 
		c = SnapConfig()
		c.read(name)
		src = c.ConfigSectionMap("main")['source']
		dest = c.ConfigSectionMap("main")['destination']
		interval = c.ConfigSectionMap("main")['interval']
		times = c.ConfigSectionMap("main")['times']
		git = c.getboolean("main", "git")
		compress = c.getboolean("main", "compress")
		srcs = [ x.strip() for x in src.split(',') ]
		# Ignore the destination if it is not a directory
		tmp_dests = [ x.strip() for x in dest.split(',') ]
		dests = [ x for x in tmp_dests if os.path.isdir(x) ]
		return [srcs, dests, interval, times, git, compress]

	def _total_size(self, source):
		total_size = os.path.getsize(source)
		for item in os.listdir(source):
		    itempath = os.path.join(source, item)
		    if os.path.isfile(itempath):
			total_size += os.path.getsize(itempath)
		    elif os.path.isdir(itempath):
			total_size += self._total_size(itempath)
		return total_size

	def add_source(self, name):
		pass

	def add_destination(self, name):
		pass

	# Support multiple destinations
	# Return a dict
	def get_free_size(self, *args):
		dict_disk_ratio = self.get_disk_ratio(args[0])
		for path in args[0]:
			# Handle links correctly
			p = os.path.realpath(os.path.abspath(path))
			dstat = os.statvfs(p)
			self.dict_free_size[path] = dstat.f_bavail * dstat.f_frsize / dict_disk_ratio[path]
		return self.dict_free_size

	def get_disk_ratio(self, *args):
		dict_disk_ratio = {}
		for path in args[0]:
			p = os.path.realpath(os.path.abspath(path))
			dict_disk_ratio[path] = os.stat(p).st_dev
		temp = { x:dict_disk_ratio.values().count(x) for x in dict_disk_ratio.values() }
		for k, v in dict_disk_ratio.items():
			if v in temp.keys():
				dict_disk_ratio[k] = temp[v]
		return dict_disk_ratio

	# Return at least one of the available destinations
	def get_available_dest(self, *args):
		src_total_size = 0
		for path in args[0]:
			if os.path.isfile(path):
				src_total_size = src_total_size + os.path.getsize(path)
			elif os.path.exists(path):
				p = os.path.realpath(os.path.abspath(path))
				src_total_size = src_total_size + self._total_size(path)
			else:
				raise ValueError
			for k, v in self.dict_free_size.items():
				if v <= src_total_size:
					del self.dict_free_size[k]
		return self.dict_free_size

	def get_timestamp(self):
		ts = time.time()
		return datetime.datetime.fromtimestamp(ts).strftime('%Y-%m-%d_%H:%M:%S')

	def start(self, **kwargs):
		def list_replace(l):
			tl = l[:]
			for i, value in enumerate(tl):
				tl[i] = value.strip('/')
				tl[i] = tl[i].replace('/', '_')
			return tl
		d = self.get_free_size(self.dest)
		available_dest = self.get_available_dest(self.src)
		st = self.timestamp 
		#Create the first backup on the first destination
		tmpdir = self.dest[0] + "/"
		if self.compress:
			tmpfile = tmpdir + "_".join(list_replace(self.src)) + "_" + st + ".tar.gz"
			tar = tarfile.open(tmpfile,"w:gz")
		else:
			tmpfile = tmpdir + "_".join(list_replace(self.src)) + "_" + st + ".tar"
			tar = tarfile.open(tmpfile, "w")
		for src in self.src:
			tar.add(src)
		tar.close()
		# Do not need to copy to the first destination
		for dest in self.dest[1:]:
			try:
				shutil.copy2(tmpfile, dest)
			except:
				pass
		self.first = False
		self.done = True

if __name__ == "__main__":
	#with daemon.DaemonContext():
	#logging.basicConfig(level=logging.INFO,
        #                format='%(asctime)s - %(message)s',
        #                datefmt='%Y-%m-%d %H:%M:%S')
	#event_handler = LoggingEventHandler()
	event_handler = SnapSetHandler()
	observer = Observer()
	for i in event_handler.s.src:
		observer.schedule(event_handler, i, recursive=True)
	observer.start()
	#s.start()
	try:
		while True:
			#time.sleep(event_handler.s.interval)
			time.sleep(1)
	except KeyboardInterrupt:
		observer.stop()
	observer.join()
