for i in `echo "$1 $2 $3 $4"`; do scp -i ~/.ssh/libra.pem $5 centos@${i}:/tmp/id; ssh -i ~/.ssh/libra.pem centos@${i} 'cat /tmp/id >> ~/.ssh/authorized_keys'; done
for i in `echo "$2 $3 $4"`; do ssh centos@${i} "sudo sed -i 's/redhat\/openshift-ovs-subnet/cni/g' /etc/origin/node/node-config.yaml"; done
for i in `echo "$1"`; do ssh centos@${i} "sudo sed -i 's/redhat\/openshift-ovs-subnet/cni/g' /etc/origin/master/master-config.yaml"; done
