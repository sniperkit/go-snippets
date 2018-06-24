#!/bin/bash

# GetKubernetes fetches k8s binaries from the k8 release repo
function GetKubernetes {

  # fetch kubernetes released binaries
  pushd .
  mkdir -p $top_dir/k8s-$k8sVer
  if [ -f $top_dir/k8s-$k8sVer/kubernetes.tar.gz ]; then
    echo "k8s-$k8sVer/kubernetes.tar.gz found, not fetching."
    rm -rf $top_dir/k8s-$k8sVer/kubernetes
    rm -rf $top_dir/k8s-$k8sVer/bin
  else
    cd $top_dir/k8s-$k8sVer 
    wget https://github.com/kubernetes/kubernetes/releases/download/$k8sVer/kubernetes.tar.gz
  fi

  # untar kubernetes released binaries
  cd $top_dir/k8s-$k8sVer
  tar xvfz kubernetes.tar.gz kubernetes/server/kubernetes-server-linux-amd64.tar.gz
  tar xvfz kubernetes/server/kubernetes-server-linux-amd64.tar.gz
  popd

  if [ ! -f $top_dir/k8s-$k8sVer/kubernetes/server/bin/kubelet ]; then
    echo "Error kubelet not found after fetch/extraction"
    exit 1
  fi
}

# GetContiv fetches k8s binaries from the contiv release repo
function GetContiv {

  # fetch contiv binaries
  pushd .
  mkdir -p $top_dir/contiv_bin
  if [ -f $top_dir/contiv_bin/netplugin-$contivVer.tar.bz2 ]; then
    echo "netplugin-$contivVer.tar.bz2 found, not fetching."
  else
    cd $top_dir/contiv_bin
    wget https://github.com/contiv/netplugin/releases/download/$contivVer/netplugin-$contivVer.tar.bz2
    tar xvfj netplugin-$contivVer.tar.bz2
  fi
  popd

  if [ ! -f $top_dir/contiv_bin/contivk8s ]; then
    echo "Error contivk8s not found after fetch/extraction"
    exit 1
  fi
}

USAGE="setup_k8s_cluster.sh <ssh_username>"
if [ $# -lt 1 ]; then
    echo $USAGE
    exit 1
fi

username=$1

# kubernetes version to use -- defaults to v1.1.4
: ${k8sVer:=v1.1.4}

# contiv version
: ${contivVer:=v0.1-05-02-2016.23-10-27.UTC}

top_dir=$PWD

#GetKubernetes
GetContiv

# generate inventory
./oc_parse_cluster.py $username

# run ansible
#ansible-playbook -i .contiv_k8s_inventory $top_dir/contrib/ansible/oc_cluster.yml --skip-tags "contiv_restart" -e "networking=contiv localBuildOutput=$top_dir/k8s-$k8sVer/kubernetes/server/bin contiv_bin_path=$top_dir/contiv_bin etcd_peers_group=masters"
ansible-playbook -i .contiv_k8s_inventory $top_dir/contrib/ansible/oc_cluster.yml -e "networking=contiv localBuildOutput=$top_dir/k8s-$k8sVer/kubernetes/server/bin contiv_bin_path=$top_dir/contiv_bin kube_master_api_port=8443 etcd_peers_group=masters"
