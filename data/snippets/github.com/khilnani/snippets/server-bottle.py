#!/usr/bin/env python

from bottle import route, run, template

@route('/hello/<name>')
def hello(name):
    return template('<b>Hello {{name}}</b>!', name=name)

@route('/')
def index():
    return ('<b>Hello from bottle.py</b>!')

run(host='0.0.0.0', port=8080, debug=True)

