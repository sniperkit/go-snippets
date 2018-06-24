#!/usr/bin/python

import sys
import json

def ipToHostname(ip):
	return "ip-"+"-".join(ip.split("."))+".ec2.internal"

class Master:

	def __init__(self, pub_ip, pvt_ip):
		self.contiv_control_if = "eth0"
		self.contiv_network_if = "tap0"
		self.contiv_network_ip = "2.2.2.2"
		self.contiv_control_ip = pvt_ip
		self.management_ip = pub_ip
		self.name = ipToHostname(pvt_ip)

	def printout(self):
		return json.dumps(self.__dict__, sort_keys=True, indent=4)

class Node:
	
	def __init__(self, pub_ip, pvt_ip, index):
		self.contiv_control_if = "eth0"
		self.contiv_network_if = "tap0"
		self.contiv_network_ip = "2.2.2.%s"%index
		self.contiv_control_ip = pvt_ip
		self.management_ip = pub_ip
		self.name = ipToHostname(pvt_ip)
		self.max_pods = 40

	def printout(self):
		return json.dumps(self.__dict__, sort_keys=True, indent=4)

def usage():
	print sys.argv[0]+ ": node1_pub_ip node1_pvt_ip node2_pub_ip node2_pvt_ip.."

if len(sys.argv) < 5:
	print "Need atleast one master and one node"
	usage()
	sys.exit(1)

inventory={"master": [Master(sys.argv[1], sys.argv[2]).__dict__]}
nodes_inventory=[]

i = 3
while i < len(sys.argv):
	nodes_inventory.append(Node(sys.argv[i], sys.argv[i+1], i).__dict__)
	i += 2

inventory["nodes"] = nodes_inventory

print json.dumps(inventory, sort_keys=True, indent=4)
