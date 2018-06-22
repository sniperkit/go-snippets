#!/usr/bin/python
import os
import sys
import json
import yaml

class SafeDict(dict):
    'Provide a default value for missing keys'
    def __missing__(self, key):
        return 'missing'

def validateHostInfo(hI):
    hostInfo = SafeDict(hI)
    for attr in hostAttr:
        if hostInfo[attr] is 'missing':
            print "{} not found for a host in cluster_defs.json".format(attr)
            sys.exit(-1)

def writeHostLine(outFd, hC, comVars):
    hConfig = SafeDict(hC)
    host = hConfig['name']
    mgmnt_ip = hConfig['management_ip']
    outFd.write("{} ansible_ssh_host={}".format(host, mgmnt_ip))
    for attr in hostAttr:
        outFd.write(" {}={}".format(attr, hConfig[attr]))

    if hConfig['max_pods'] is not 'missing':
        outFd.write(" {}=\"{}\"".format('max_pods', hConfig['max_pods']))

    outFd.write(comVars)
    outFd.write("\n")

def readConfig():
    cFd = open("cluster_defs.json")
    res = json.load(cFd)
    cFd.close()
    return res

def parseACI():
    try:
        aciFd = open("aci.yml")
    except Exception:
        return ""
    
    aciInfo = yaml.load(aciFd)
    aciFd.close()

    if type(aciInfo) is not dict:
        print "Warning: aci.yml was ignored"
        return ""

    res = ""
    for key, val in aciInfo.iteritems():
       if type(val) is str: 
           res += " {}={}".format(key.lower(), val)
       elif type(val) is list: 
           mList = ""
           mCount = 0
           for m in val:
               if type(m) is str:
                   if mCount > 0:
                       mList += ","
                   mList += m
                   mCount += 1
               else:
                   print "aci.yml: Error parsing key {} unexepected type {} ".format(key, type(m))
                   sys.exit(1)

           res += " {}={}".format(key.lower(), mList)
           
       else:
           print "aci.yml: Error parsing key {} unexepected type {} ".format(key, type(val))
           sys.exit(1)

    return res

hostAttr = ["name", "management_ip", "contiv_control_if", "contiv_control_ip", "contiv_network_if", "contiv_network_ip"]

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print "Usage: parse_cluster.py <user_id>"
        sys.exit(1)

    ssh_user = sys.argv[1]

    # read aci info
    aciVars = parseACI()

    # read cluster config
    clusterConf = readConfig()

    # Make sure all hosts are reported
    for mInfo in clusterConf['master']:
        validateHostInfo(mInfo)
        service_ip = mInfo['contiv_control_ip']

    for nInfo in clusterConf['nodes']:
        validateHostInfo(nInfo)

    common_vars = " ansible_ssh_user=" + ssh_user
    common_vars += " contiv_service_vip={}".format(service_ip)
    # add proxy if applicable
    proxy = os.environ.get('http_proxy')
    if proxy is not None:
        common_vars += " http_proxy={} https_proxy={}".format(proxy, proxy)
    else:
        proxy = os.environ.get('https_proxy')
        if proxy is not None:
            common_vars += " https_proxy={}".format(proxy)

    outFd = open(".contiv_k8s_inventory", "w")
    outFd.write("[masters]\n")
    for mInfo in clusterConf['master']:
        # include only management and control ip in the certificate
        cert_ip = " kube_cert_ip=" + mInfo['management_ip']
        if mInfo['contiv_control_ip'] != mInfo['management_ip']:
            cert_ip += ",IP:" + mInfo['contiv_control_ip']

        writeHostLine(outFd, mInfo, common_vars + cert_ip + aciVars)

    outFd.write("[etcd_servers]\n")
    for mInfo in clusterConf['master']:
        writeHostLine(outFd, mInfo, common_vars)

    outFd.write("[etcd]\n")
    for mInfo in clusterConf['master']:
        writeHostLine(outFd, mInfo, common_vars + " etcd_proxy_mode=off")

    outFd.write("[nodes]\n")
    for nInfo in clusterConf['nodes']:
        writeHostLine(outFd, nInfo, common_vars)

    outFd.close()

    etcFd = open(".etc_hosts", "w")
    for mInfo in clusterConf['master']:
        etcFd.write("{}       {}\n".format(mInfo['contiv_control_ip'], mInfo['name']))
    for nInfo in clusterConf['nodes']:
        etcFd.write("{}       {}\n".format(nInfo['contiv_control_ip'], nInfo['name']))

    etcFd.close()
