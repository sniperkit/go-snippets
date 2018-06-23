from flask import Flask, url_for
from flask import render_template
from random import randint
import json

app = Flask(__name__)
app._static_folder = "/var/www/junyi/flask"
max = 10

@app.route('/')
def hello_world():
	return render_template('index.html')

@app.route('/addition')
@app.route('/addition/<name>')
def addition(name="Junyi"):
	num_list = []
	while len(num_list) < 3:
		num_list.append([randint(1, 5), randint(1,5)])
		
	return render_template('addition.html', name=name, num_list=num_list)

if __name__ == '__main__':
	app.run(host='0.0.0.0', debug=True)
