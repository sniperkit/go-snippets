#!flask/bin/python

from flask import Flask, jsonify, abort
import subprocess

app = Flask(__name__)

tasks = [
{
	'id': 1,
	'name': "Hao",
	'sex': "Male"
},
{
	'id': 2,
	'name': "Xue",
	'sex': "Female"
}]

@app.route('/')
def index():
	return "Hello, World!"

@app.route('/tasks', methods=['GET'])
def get_tasks():
	return jsonify({"tasks": tasks})
@app.route('/task/<int:task_id>', methods=['GET'])
def get_task(task_id):
	task = filter(lambda x: x['id']==task_id, tasks)
	if(len(task)==0):
		abort(404)
	return jsonify({"task": task})

@app.route('/cpu')
def get_cpu():
	output = subprocess.Popen(["cat", "/proc/cpuinfo"], stdout=subprocess.PIPE)
	return "".join(output.stdout.readlines())

@app.route('/mem')
def get_mem():
	output = subprocess.Popen("free", shell=True, stdout=subprocess.PIPE)
	return "".join(output.stdout.readlines())

@app.route('/disk')
def get_disk():
	output = subprocess.Popen("df -h", shell=True, stdout=subprocess.PIPE)
	return "".join(output.stdout.readlines())

if __name__ == '__main__':
	app.run(debug=True)
