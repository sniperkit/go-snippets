---
title: "A Night with Kubernetes"
date: 2018-03-21T23:13:40+08:00
tags: ["devops", "cloud"]
draft: false
---

### What is Kubernetes?

![kubernetes logo](/images/kubernetes_logo-sm.png#featured)
<p class="figure-text"><i>*Kubernetes Logo.</i></p>

Dealing with lot of applications that are containerized (ie. via docker) and installed across regions can cause headache; from restarting containers, deploying to multiple nodes, scaling, etc. Kubernetes is a tool to help us deal with such situation by putting those containers in a _cluster_.

<!--more-->

### Kubernetes Cluster, Nodes, Deployments, and Pods

Knowing these four terms will help us understand how Kubernetes work. **Cluster** is basically a pool of machines, which at least one of it acts as a cluster master, and the other machines act as worker machines (slaves). As you may've guessed, master machine will manage the operation of its slaves. **Node** is a term used to address a worker machine in a cluster. **Deployments** refer to the applications we want to deploy or have deployed. **Pod**, the smallest unit in Kubernetes, is a group of one or more containers that makes up an application or service. You can think of pod as your logical host, or virtual machine. The overall logical view of a Kubernetes cluster would look like this:

```
▾ Kubernetes Cluster
  ▾ Master #1
    ▾ Deployment "app1"
      ▾ Node #1
        ▾ Pod "app1-2035384211"
          nginx
          redis
        Pod ...
    ▾ Deployment "app2"
      ▾ Node #1
        Pod "app2-2035384224"
          [redacted]
        Pod ...
    ▾ Deployment "app3"
      ▾ Node #1
        Pod "app3-2035384277"
          [redacted]
        ...
      ▾ Node #2
        Pod "app3-2035384122"
          [redacted]
    ▾ Deployment "app-xyz"
      ...
  ▾ Master #N
    ...
```
<p class="figure-text">Figure 1: Logical view of typical Kubernetes cluster</p>

As shown in above structure, it's a a Kubernetes cluster with single master. However, we still can add more master (master N). Our master #1 currently contains 3 deployments; `app1`, `app2`, and `app3`. In practice, you'd see those three apps as web server, cache, or your custom application.

The `app1` has been deployed to one node, `Node #1`. Each deployment will have a unique suffix in its pod name, just like `app1-2035384211`. Digging deeper into pod `app1-2035384211`, we know that it runs two containers; nginx and redis. For every pod we have in our Kubernetes cluster, we can execute command (bash, etc) in it using `kubectl`. Below command will run bash in nginx container of pod `app1-2035384211`: 

`kubectl exec -it app1-2035384211 --container nginx -- /bin/bash`

Going to deployment `app2`, it shows that pod `app2-2035384224` has been deployed to `Node #1`, which means that at this point, our `Node #1` contains 2 deployments; `app1` and `app2`. Deployment `app3` is a bit different than the previous two deployments because it has been deployed to two nodes (`Node #1` and `Node #2`). In other words, our `app3` has been scaled into two instances, which spread into two nodes. Together with `app3`, our `Node #1` now contains 3 deployments; `app1`, `app2`, and `app3`. Our `Node #2` however, only contains single pod so far, which is `app3-2035384122`.

### Creating Kubernetes Cluster

To simulate the structure in figure 1, we'll create a Kubernetes Cluster of three machines. One machine for cluster master, and other two as nodes. Let's name these machines as `vm1`, `vm2`, and `vm3` for this demonstration. Next, creating the cluster is quite straightforward:
    
**1. Install kubernetes dependencies**

```
# taken from 
# https://zihao.me/post/creating-a-kubernetes-cluster-from-scratch-with-kubeadm/

# 1. prepare new repos
apt-get -y install apt-transport-https ca-certificates software-properties-common curl

# 2. add docker repo
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository \
       "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
       $(lsb_release -cs) \
       stable"

# 3. add kubernetes repo
curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key  add -
add-apt-repository \
      "deb https://apt.kubernetes.io/ \
      kubernetes-$(lsb_release -cs) \
      main"

# install dependencies
apt-get update && apt-get install -y docker-ce kubelet kubeadm kubernetes-cni
```

If everything went well, the `docker`, `kubectl`, and `kubeadm` commands should be available in your machines.

**2. Initialize `vm1` as master**

- ssh to `vm1`
- run `kubeadm init`

The `kubeadm init` command will output something like this: 

```
kubeadm join --token token.abc.def {server_ip}:6443 --discovery-token-ca-cert-hash sha256:blabla
```

**3. Make the other machines join the cluster**

- ssh to `vm2`
- run the `kubeadm join` command from step 2
- install container network interface (CNI[^2])

    ```
    # using weave CNI for now, but you can try other choices
    kubectl apply -f https://git.io/weave-kube-1.6
    ```

- repeat for `vm3`

<p class="text-center">***</p>

*The two important parts to create Kubernetes cluster from above steps is that we did `kubeadm init` in `vm1` to make it as cluster master, and then `kubeadm join ...` from other machines to join the cluster.*

Now that we've a cluster master and joined the other two machines to our cluster master, we can verify if it's reflected in our list nodes by executing this command in cluster master: `kubectl get nodes`

![kubernetes nodes](/images/kubernetes-nodes.jpg)
<p class="figure-text">Figure 2: Displaying Kubernetes nodes</p>

> If you forgot to install CNI in step 3, your node status will likely yield `Pending` instead of `Ready`

In figure 2, you must have noticed that before we run `kubectl get nodes`, we set an environment variable `KUBECONFIG`. This variable refers to the config file used by `kubectl`, which defaults to `~/.kube/config`. However, we changed it since the we want the `kubectl` to interact with our cluster master, whose config is stored `/etc/kubernetes/admin.conf` (generated when we run `kubeadm init` command). Another approach if we don't want to change the `KUBECONFIG` is to pass a config file to kubectl's `--kubeconfig` arg, as follow:

`kubectl get nodes --kubeconfig=/etc/kubernetes/admin.conf`

### Creating Deployment

You can create Kubernetes deployment using file or directly passing container's image url. For this demonstration, we'll create a deployment named `app1` which simply is a hello-world-nginx app pulled from public docker hub at `kitematic/hello-world-nginx`. To do so, we just need to execute this command:

`kubectl run app1 --image=kitematic/hello-world-nginx --port=80`

Above command will create one deployment named `app1` and one pod which name prefix is `app1`. The `--port=80` argument indicates that we want to expose port 80 of the image. To verify if the deployment and pod exist, run this command:

```
kubectl get deployments
kubectl get pods
```

![kubernetes deployments](/images/kubernetes-deployments.jpg)
<p class="figure-text">Figure 3: Listing Kubernetes deployments and pods</p>

As shown in above image, we've successfully deployed our `app1` and the pods is running in one of our node. To get more detailed info of our `app1-ff7b44c5d-q8rxd` pod (the node where the pod has been deployed, state, IP, restart count, etc.), run this command:

`kubectl describe pods app1-ff7b44c5d-q8rxd`

### Exposing Deployments

Creating deployments will only create pods and assign it to internal IP of the cluster, which doesn't make the app accessible from outside vm; we can't access app1 from `{ip-of-vm}:{port-of-app1}`. To make the app publicly accessible, we'll need to expose our deployment. For this demonstration, we'll use the `NodePort` type which will bind our deployment port to vm's port:

`kubectl expose deployment app1 --type=NodePort --port 80`

Once `app1` is exposed, we can verify by checking the list of ports have been bound to our vm (the Kubernetes Services[^3]), using `kubectl get service`:

![kubernetes services](/images/kubernetes-expose-service.jpg)
<p class="figure-text">Figure 4: List of Kubernetes Services</p>

We can see that port 80 of our `app1` is bound to port `30485` (of our vm1), which means that the `app1` is now publicly accessible from `{ip-of-vm1}:30485`. If the IP of your `vm1` is `67.207.85.43`, then you can access `app1` from `67.207.85.43:30485`.

### Closing

That's it! We've explored quite a lot of Kubernetes; the structure, nodes, deployments, and all the way down
 to make the app accessible from internet. Of course, accessing your app from `67.207.85.43:30485` instead of `67.207.85.43` is not intuitive, which is why, in practice, we'd put another layer to handle it (ie. proxy, load balancer, or ingress[^4]).

The next thing that you may want to explore is [Kubernetes Dashboard](https://github.com/kubernetes/dashboard), which enables us to control all Kubernetes things visually; deploying nodes, executing command in pods, deleting deployments, adjusting number of replicas, etc. 

Lastly, while this post is not meant to be an in-depth notes of Kubernetes, hopefully it gives you a better idea of what Kubernetes is and what is it capable of. Also, they've a very complete documentation, which is awesome! See it here: https://kubernetes.io/docs/home/

_"It turns out to be a long night of exploring Kubernetes"_ 😆

***Till next. See ya!***

[^1]:https://zihao.me/post/creating-a-kubernetes-cluster-from-scratch-with-kubeadm/
[^2]:https://github.com/containernetworking/cni
[^3]:https://kubernetes.io/docs/concepts/services-networking/service/
[^4]:https://kubernetes.io/docs/concepts/services-networking/ingress/