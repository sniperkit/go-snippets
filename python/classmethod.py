class a():
	@staticmethod
	def staticm():
		print 'static'
	def normalm(self):
		print 'normal', self
	@classmethod
	def classm(cls):
		print 'class', cls

class A(object):
	@classmethod
	def cm(cls):
		print "caller: ", cls.__name__
	@staticmethod
	def sm():
		print "static"

class B(A):
	pass

a1 = a()
a1.normalm()
a1.staticm()
a1.classm()
print type(a)
print type(a1)

A.cm()
B.cm()
A.sm()
B.sm()

class A:
    member = "this is a test."
    def __init__(self):
        pass

    @classmethod
    def Print1(cls):
        print "print 1: ", cls.member
        
    def Print2(self):
        print "print 2: ", self.member
           
        
    @classmethod    
    def Print3(paraTest):
        print "print 3: ", paraTest.member
    @staticmethod
    def print4():
        print "hello"
    

a = A()
A.Print1()  
a.Print1()
#A.Print2() Error expected
a.Print2()
A.Print3()
a.Print3() 
A.print4()
a.print4()
