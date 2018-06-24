!/bin/bash

echo "####Creating Containers####"
docker run -dt  --cap-add=ALL  --name ospf1 --ip 192.168.0.2 --net=clos-oob-network  -P ubuntu_flex:v2
docker run -dt  --cap-add=ALL --name ospf2 --ip 192.168.0.3 --net=clos-oob-network  -P ubuntu_flex:v2

ospf1_pid=`docker inspect -f '{{.State.Pid}}' ospf1`
ospf2_pid=`docker inspect -f '{{.State.Pid}}' ospf2`

mkdir -p /var/run/netns

ln -s /proc/$ospf1_pid/ns/net /var/run/netns/$ospf1_pid
ln -s /proc/$ospf2_pid/ns/net /var/run/netns/$ospf2_pid

echo -e "done!\n"

# Create the "peer" interfaces and hand them out

echo "###Creating P2P interfaces###"

#ospf1 to ospf3-4 ports 8
sudo ip link add eth1 type veth peer name eth11
sudo ip link add eth2 type veth peer name eth12
sudo ip link add eth3 type veth peer name eth13

#ospf1 interfaces
sudo ip link set eth1 netns $ospf1_pid
sudo ip netns exec $ospf1_pid ip link set eth1 up
sudo ip link set eth2 netns $ospf1_pid
sudo ip netns exec $ospf1_pid ip link set eth2 up
sudo ip link set eth3 netns $ospf1_pid
sudo ip netns exec $ospf1_pid ip link set eth3 up


#ospf2 interfaces
sudo ip link set eth11 netns $ospf2_pid
sudo ip netns exec $ospf2_pid ip link set eth11 up
sudo ip link set eth12 netns $ospf2_pid
sudo ip netns exec $ospf2_pid ip link set eth12 up
sudo ip link set eth13 netns $ospf2_pid
sudo ip netns exec $ospf2_pid ip link set eth13 up


echo -e "done!\n"

