# Deploying to Kubernetes

To deploy `pushprom` to Kubernetes you need two things:

* a [NodePort Service](http://kubernetes.io/docs/user-guide/services/#type-nodeport)
* a [Deployment](http://kubernetes.io/docs/user-guide/deployments/)

The Service will expose UDP port 9090 and TCP port 9091 for internal usage, while also generating NodePorts for both (see # for usage). The Deployment will ensure a single [Pod](http://kubernetes.io/docs/user-guide/pods/) with `pushprom` is always running. Note that, since `pushprom` is not a distributed service, you can only have one replica of `pushprom` running in a Pod.

## Deploying `pushprom`

### Creating a namespace
Before you get started, ensure that the Kubernetes namespace you want to run on exists, you can do that by running:
```
kubectl get namespace <NAMESPACE>
```
If it does not exist, create it with the following command:
```
kubectl create namespace <NAMESPACE>
```

### Creating Service and the Deployment
Create two files, one called [service.yaml](service.yaml) and one called [deployment.yaml](deployment.yaml). These files will contain the specification for the Service and the Deployment respectively.

To deploy to them to your Kubernetes cluster you run the following commands:
```
kubectl apply -f deployment.yaml
```
```
kubectl apply -f service.yaml
```

To verify that both the Deployment and the Service generation (or update) worked, run:
```
kubectl get po,svc --namespace=<NAMESPACE>
```

## Making `pushprom` available to your service
By creating a Service you've now made `pushprom` available inside your Kubernetes cluster. You can talk to it by using the following address: `pushprom.<NAMESPACE>`.

## Making `pushprom` available to Prometheus
In order for [Prometheus](https://prometheus.io/) to reach your `pushprom` installation, you need to expose the Service to the outside world. First, let's find out what NodePorts were created by running:
```
kubectl describe svc pushprom --namespace=<NAMESPACE>
```
Using the NodePorts that are listed there, create a rule in your load balancer (or similar ingress point) to point to these ports.
