# Docker Swarm Cluster State Manager

This is not an official Google product.

## Overview

The Swarm cluster state manager is a simple prototype that demonstrates how cluster state could be managed for Docker Swarm. The ideas for this prototype are based on the Kubernetes pod replication controller.

## Usage

The following tutorial assumes you have [created a Docker Swarm cluster using Docker Machine](https://docs.docker.com/swarm/install-w-machine).

### Start the swarm cluster state manager

Edit the `docker-compose.yml` file and replace `SWARM_MASTER_IP` with the swarm-master IP address.

```
$ sed -i -e "s/SWARM_MASTER_IP/$(docker-machine ip swarm-master)/g" docker-compose.yml
```

Launch the `swarm-cluster-state-manager` service on the same Docker host as the `swarm-master` service.

Use the `eval` command to ensure your are pointing to the same Docker host where the `swarm-master` service
is running:

```
$ eval $(docker-machine env swarm-master)
```

Use the `docker-compose up` command to start the `swarm-cluster-state-manager` service:

```
$ docker-compose up -d
```

### Reusing the Docker machine client certs

#### OS X

On OS X you'll need to create a pkcs12 bundle of the PEM encoded client certs in order
to use them with cURL.

```
$ openssl pkcs12 -export \
  -in ~/.docker/machine/certs/cert.pem \
  -inkey ~/.docker/machine/certs/key.pem \
  -out ~/.docker/machine/certs/cert.pfx 
```

At the follow prompt set the export password to protect the cert and key:

```
Enter Export Password:
Verifying - Enter Export Password:
```

Later examples assume the word `swarm` was used for the password.

### Submit a new cluster state object

Use the `eval` command to ensure you are pointing to the swarm-master endpoint:

```
$ eval $(docker-machine env swarm-master --swarm)
```

The following command submits a cluster state object named nginx and will ensure 5
containers are started from the `nginx:1.9.6` Docker image.

```
$ curl https://$(docker-machine ip swarm-master):3476/submit \
  -d '{"Name": "nginx", "Image": "nginx:1.9.6", "Count": 5}' \
  --cacert ~/.docker/machine/certs/ca.pem \
  --cert ~/.docker/machine/certs/cert.pfx \
  --pass swarm
```

Review the logs for the `swarm-cluster-state-manager` service:

```
$ docker logs swarm-cluster-state-manager
```

```
Starting Swarm cluster state manager...
2015/11/19 17:55:46 updating nginx cluster state
2015/11/19 17:55:50 creating 5 nginx containers
2015/11/19 17:55:50 created nginx container: f7ca99588905796b9a0fd6dcaf057e433170e876fd9500352770bd8d8096cc79
2015/11/19 17:55:51 created nginx container: 1406e27db94b792906e4c86e76cf9c4ba703b3aefc4ee651ecc10826cd18c612
2015/11/19 17:55:51 created nginx container: 7307b5221fb30f3373ecda4a9595774fcb8e673782b8ab73a4bc6053efb15749
2015/11/19 17:55:51 created nginx container: c17f8163e375736a30c16a2b5ca58b70e5330908ae238d0f43e8c99049622514
2015/11/19 17:55:51 created nginx container: 55b49c75914153603998e9e40687424d8a46890d392abd75b0f6276dc36c6b23
```

### Automatically replace deleted containers

Use the `eval` command to ensure you are pointing to the swarm-master endpoint:

```
$ eval $(docker-machine env swarm-master --swarm)
```

All containers created by the `swarm-cluster-state-manager` service include a
`swarm.cluster.state` label with the service name as the value.

Run the following command to delete all containers defined by the nginx cluster
state entry:

```
$ docker rm -f $(docker ps -f label=swarm.cluster.state=nginx -q)
```

Observe the `swarm-cluster-state-manager` service logs

```
$ docker logs swarm-cluster-state-manager
```

```
Starting Swarm cluster state manager...
2015/11/19 17:55:46 updating nginx cluster state
2015/11/19 17:55:50 creating 5 nginx containers
2015/11/19 17:55:50 created nginx container: f7ca99588905796b9a0fd6dcaf057e433170e876fd9500352770bd8d8096cc79
....
2015/11/19 18:00:01 creating 5 nginx containers
2015/11/19 18:00:01 created nginx container: 67532a34aa04f346d1d0cae7947c4f64361da9058ba0b84ec4d4b57f9b443337
2015/11/19 18:00:01 created nginx container: b3204f60b3d0130d357bcd429a46af54cb640e9068abed4ef281b9eba7969761
2015/11/19 18:00:02 created nginx container: 88e00433c5c0e72c280f0af8ee9b5e1c73f1609d93d69665a8e175a4f958c67b
2015/11/19 18:00:02 created nginx container: ad6ac5a04167bd5f81f61f547c91403ec93f3869c8174960d49a732959092b29
2015/11/19 18:00:02 created nginx container: 7cc9afea05027875fc4100b30fe0fd65ef2bad9a68012c58b65b581ccddf4805
```

### Update the Cluster state

Reduce the number of desired running nginx containers to 3: 

```
$ curl https://$(docker-machine ip swarm-master):3476/submit \
  -d '{"Name": "nginx", "Image": "nginx:1.9.6", "Count": 3}' \
  --cacert ~/.docker/machine/certs/ca.pem \
  --cert ~/.docker/machine/certs/cert.pfx \
  --pass swarm
```

Observe the swarm-cluster-state-manager service logs

```
$ docker logs swarm-cluster-state-manager
```

```
...
2015/11/19 18:03:05 updating nginx cluster state
2015/11/19 18:03:12 removing 2 nginx containers
2015/11/19 18:03:12 removed nginx container: b3204f60b3d0130d357bcd429a46af54cb640e9068abed4ef281b9eba7969761
2015/11/19 18:03:12 removed nginx container: 67532a34aa04f346d1d0cae7947c4f64361da9058ba0b84ec4d4b57f9b443337
```

### Get the current status

The following command retrieves the cluster status from the Swarm cluster state manager.

```
$ curl https://$(docker-machine ip swarm-master):3476/status \
  --cacert ~/.docker/machine/certs/ca.pem \
  --cert ~/.docker/machine/certs/cert.pfx \
  --pass swarm
```

```
{
...
    "CurrentCount": 5,
    "DesiredCount": 5,
    "Image": "nginx:1.9.6",
    "Name": "nginx"
}
...
```
