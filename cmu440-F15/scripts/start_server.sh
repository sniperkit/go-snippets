#!/bin/bash

if [ -z $GOPATH ]; then
    echo "FAIL: GOPATH environment variable is not set"
    exit 1
fi

if [ -n "$(go version | grep 'darwin/amd64')" ]; then    
    GOOS="darwin_amd64"
elif [ -n "$(go version | grep 'linux/amd64')" ]; then
    GOOS="linux_amd64"
else
    echo "FAIL: only 64-bit Mac OS X and Linux operating systems are supported"
    exit 1
fi

# Build srunner
go install github.com/cmu440-F15/paxosapp/runners/srunner
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

# Build the student's paxos node implementation.
# Exit immediately if there was a compile-time error.
go install github.com/cmu440-F15/paxosapp/runners/prunner
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

# Build mrunner
go install github.com/cmu440-F15/paxosapp/runners/mrunner
if [ $? -ne 0 ]; then
echo "FAIL: code does not compile"
exit $?
fi

# Build the student's paxos node implementation.
# Exit immediately if there was a compile-time error.
go install github.com/cmu440-F15/paxosapp/runners/mrunner
if [ $? -ne 0 ]; then
echo "FAIL: code does not compile"
exit $?
fi

# Pick random ports between [10000, 20000) for Monitor node
MONITOR_PORT=$(((RANDOM % 10000) + 10000))

# Pick random ports between [10000, 20000).
SLAVE_PORT0=$(((RANDOM % 10000) + 10000))
SLAVE_PORT1=$(((RANDOM % 10000) + 10000))
SLAVE_PORT2=$(((RANDOM % 10000) + 10000))
SLAVE_PORT3=$(((RANDOM % 10000) + 10000))
SLAVE_PORT4=$(((RANDOM % 10000) + 10000))
SLAVE_PORT5=$(((RANDOM % 10000) + 10000))
SLAVE_NODE=$GOPATH/bin/srunner
ALL_SLAVE_PORTS="${SLAVE_PORT0},${SLAVE_PORT1},${SLAVE_PORT2},${SLAVE_PORT3},${SLAVE_PORT4},${SLAVE_PORT5}"

echo $ALL_SLAVE_PORTS
# Start slave node.
${SLAVE_NODE} -hostport="localhost:$SLAVE_PORT0" -id=0 2 &
SLAVE_NODE_PID0=$!
sleep 1

${SLAVE_NODE} -hostport="localhost:$SLAVE_PORT1" -id=1 2 &
SLAVE_NODE_PID1=$!
sleep 1

${SLAVE_NODE} -hostport="localhost:$SLAVE_PORT2" -id=2 2 &
SLAVE_NODE_PID2=$!
sleep 1

${SLAVE_NODE} -hostport="localhost:$SLAVE_PORT3" -id=3 2 &
SLAVE_NODE_PID3=$!
sleep 1

${SLAVE_NODE} -hostport="localhost:$SLAVE_PORT4" -id=4 2 &
SLAVE_NODE_PID4=$!
sleep 1

${SLAVE_NODE} -hostport="localhost:$SLAVE_PORT5" -id=5 2 &
SLAVE_NODE_PID5=$!
sleep 1


# Pick random ports between [10000, 20000).
NODE_PORT0=$(((RANDOM % 10000) + 10000))
NODE_PORT1=$(((RANDOM % 10000) + 10000))
NODE_PORT2=$(((RANDOM % 10000) + 10000))

echo "localhost:$NODE_PORT0,localhost:$NODE_PORT1,localhost:$NODE_PORT2" > paxosnodesport.txt
#TESTER_PORT=$(((RANDOM % 10000) + 10000))
#PROXY_PORT=$(((RANDOM & 10000) + 10000))
#PAXOS_TEST=$GOPATH/bin/paxostest
PAXOS_NODE=$GOPATH/bin/prunner
ALL_PORTS="${NODE_PORT0},${NODE_PORT1},${NODE_PORT2}"

##################################################
echo "${PAXOS_NODE} -ports=${ALL_PORTS} -slaveports=${ALL_SLAVE_PORTS} -N=3 -id=0 -numslaves=6 2 &"

# Start paxos node.
${PAXOS_NODE} -ports=${ALL_PORTS} -slaveports=${ALL_SLAVE_PORTS} -monitorport=${MONITOR_PORT} -N=3 -id=0 -numslaves=6 2 &
PAXOS_NODE_PID0=$!
sleep 1

${PAXOS_NODE} -ports=${ALL_PORTS} -slaveports=${ALL_SLAVE_PORTS} -monitorport=${MONITOR_PORT} -N=3 -id=1 -numslaves=6 2 &
PAXOS_NODE_PID1=$!
sleep 1

${PAXOS_NODE} -ports=${ALL_PORTS} -slaveports=${ALL_SLAVE_PORTS} -monitorport=${MONITOR_PORT} -N=3 -id=2 -numslaves=6 2 &
PAXOS_NODE_PID2=$!
sleep 1


rm server_log.txt 2> /dev/null
# slave node to kill.
echo ${SLAVE_NODE_PID0} >> server_log.txt
echo ${SLAVE_NODE_PID1} >> server_log.txt
echo ${SLAVE_NODE_PID2} >> server_log.txt
echo ${SLAVE_NODE_PID3} >> server_log.txt
echo ${SLAVE_NODE_PID4} >> server_log.txt
echo ${SLAVE_NODE_PID5} >> server_log.txt


# paxos node to kill.
echo ${PAXOS_NODE_PID0} >> server_log.txt
echo ${PAXOS_NODE_PID1} >> server_log.txt
echo ${PAXOS_NODE_PID2} >> server_log.txt




MONITOR_NODE=$GOPATH/bin/mrunner

$MONITOR_NODE -port=${MONITOR_PORT}  -masterPort=${ALL_PORTS} &

MONITOR_PID=$!
echo ${MONITOR_PID} >> server_log.txt

while [[ true ]]; do
  while read line
  do
    ${PAXOS_NODE} -ports=${ALL_PORTS} -slaveports=${ALL_SLAVE_PORTS} -N=3 -id=$line -monitorport=${MONITOR_PORT} -numslaves=6 -replace=true 2 &
    NEW_PAXOS_NODE_PID=$!
    echo ${NEW_PAXOS_NODE_PID} >> server_log.txt
  done < server_id.txt 2> /dev/null
  rm server_id.txt 2> /dev/null
  sleep 2
done 2> /dev/null 

# /Users/Yifan/Documents/gowork/bin/prunner -ports=14987,15870,11164 -slaveports=13902,13029,12353,14254,16935,19558 -N=3 -id=0 -numslaves=6 -replace=true


