import webob.exc
import webob.dec
from webob import Request
import paste.httpserver

class TestRequest(webob.Request):
	#@property
	def is_local(self):
		return self.remote_addr == '127.0.0.1'

@webob.dec.wsgify(RequestClass=TestRequest)
def testfunc(req):
	import pdb; pdb.set_trace()
	if req.is_local():
		return Response('hi!')
	else:
		raise webob.exc.HTTPForbidden

req = Request.blank('test')
req.remote_addr = '127.0.0.1'
res = testfunc(req)
