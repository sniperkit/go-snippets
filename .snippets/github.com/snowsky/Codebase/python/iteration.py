class MyIterator(object):
	def __init__(self, step):
		self.step = step
	def __iter__(self):
		return self
	def next(self):
		if self.step == 0:
			raise StopIteration
		self.step -= 1
		return self.step
