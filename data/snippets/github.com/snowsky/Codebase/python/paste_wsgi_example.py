'''
A minimal example of how to use Paste and WebOb to build a custom
WSGI app and serve it.
 
Depends on:
* paste - http://pypi.python.org/pypi/Paste
* webob - http://pypi.python.org/pypi/WebOb/1.1.1
* routes - http://pypi.python.org/pypi/Routes/1.12.3
 
I (marmida) still think this is less appropriate than using CouchDB; you'll need
to handle routing and controllers manually to reproduce Couch's default behavior 
(a RESTful interface to db records).
 
I extracted this from an app I was working on a while ago:
https://github.com/marmida/catalog/blob/master/server/app.py
 
I originally built that code to work with neo4j, and then switched to sqlite.
At the time, neo4j required thread isolation; I kept that after I moved to sqlite.
If you find yourself hitting threading issues (e.g. sqlite doesn't like being 
called from paste worker threads), you can look at my code on GitHub for a solution.
'''
 
import paste.fileapp
import paste.httpserver
import routes
import webob
import webob.dec
import webob.exc
 
 
HOST = '0.0.0.0'
PORT = 8080
 
class ExampleApp(object):
    '''
    A WSGI "application."
 
    See: http://pythonpaste.org/do-it-yourself-framework.html#writing-a-wsgi-application
    '''
 
    # Our routes map URIs to methods of this app, and define how to extract args from requests
    # Complaint: in order to make this RESTful, you have to plan routes yourself
    map = routes.Mapper()
    map.connect('index', '/', method='index')
    map.connect('greet', '/greet/{name}', method='greet')
    
    #@webob.dec.wsgify
    def __call__(self, req):
        '''
        Glue.  A WSGI app is a callable; thus in order to make this object an application, 
        we define __call__ to make it callable.  We then ask our Mapper to do some routing,
        and dispatch to the appropriate method.  That method must return a webob.Response.
        '''
	import pdb; pdb.set_trace()
        results = self.map.routematch(environ=req.environ)
        if not results:
            return webob.exc.HTTPNotFound()
        match, route = results
        link = routes.URLGenerator(self.map, req.environ)
        req.urlvars = ((), match)
        kwargs = match.copy()
        method = kwargs.pop('method')
        req.link = link
        return getattr(self, method)(req, **kwargs)
 
    def index(self, req):
        '''
        Controller #1: a landing page.
        '''
        return webob.Response(
            body='''<html><body>
                    <p>Your name: 
                        <form onSubmit="location.href='/greet/' + encodeURI(document.getElementById('name_input').value); return false;">
                            <input type="text" value="" id="name_input"/>
                        </form>
                    </p>
                </body></html>'''
        )
        
    def greet(self, req, name=None):
        '''
        Controller #2: do something with a URI arg to show dynamic behavior.
        '''
        return webob.Response(
            body='<html><body>Dear %s, you\'re a shmuck.</body></html>' % name
        )
 
def main():
    '''
    CLI entry point.
    '''
    paste.httpserver.serve(ExampleApp(), host=HOST, port=PORT)
 
# for importability, don't run main() unless we really mean to
if __name__ == '__main__':
    main()
