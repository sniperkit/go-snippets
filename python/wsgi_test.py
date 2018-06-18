from wsgiref.simple_server import make_server

def my_app(environ, start_response):
    """a simple wsgi application"""
    status = '200 OK'
    response_headers = [('Content-type', 'text/plain')]
    start_response(status, response_headers)
    return ["hello world"]

httpd = make_server('', 8000, my_app)
print "Serving on port 8000..."
httpd.serve_forever()
