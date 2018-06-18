import ConfigParser

class SnapConfig(ConfigParser.ConfigParser):
	def __init__(self):
		#super(SnapConfig, self).__init__()
		ConfigParser.ConfigParser.__init__(self)
	def ConfigSectionMap(self, section):
	    dict1 = {}
	    options = self.options(section)
	    for option in options:
		try:
		    dict1[option] = self.get(section, option)
		    if dict1[option] == -1:
			DebugPrint("skip: %s" % option)
		except:
		    print("exception on %s!" % option)
		    dict1[option] = None
	    return dict1

if __name__ == '__main__':
	Config = SnapConfig()
	Config.read("./test.ini")
	print Config.sections()
	print Config.ConfigSectionMap("SectionOne")['name']
	dests = Config.ConfigSectionMap("SectionOne")['destination']
	test = [ x.strip() for x in dests.split(',') ]
	print test
	for dest in dests.split(','):
		print dest.strip()
