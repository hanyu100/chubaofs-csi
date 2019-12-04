# Overview

ChubaoFS Container Storage Interface (CSI) plugins.

## Prerequisite

* Kubernetes 1.12.5
* ChubaoFS 1.4.0
* CSI spec version 0.3.0.

## Enable privileged Pods

To use CSI drivers, your Kubernetes cluster must allow privileged pods (i.e. --allow-privileged flag must be set to true for both the API server and the kubelet). Ensure your API server are started with the privileged flag.

```
$ ./kube-apiserver ...  --allow-privileged=true ...
$ ./kubelet ...  --allow-privileged=true ...
```

## Prepare on-premise ChubaoFS cluster

An on-premise ChubaoFS cluster can be deployed separately, or within the same Kubernetes cluster as applications which requrie persistent volumes. Please refer to [chubaofs-helm](https://github.com/chubaofs/chubaofs-helm) for more details on deployment using Helm.

## Deploy the CSI driver

```
$ kubectl create configmap kubecfg --from-file=deploy/kubecfg```
$ kubectl apply -f deploy/cfs-rbac.yaml
$ kubectl apply -f deploy/csi-controller-deployment.yaml
$ kubectl apply -f deploy/csi-node-daemonset.yaml
```
## Use Remote ChubaoFS Cluster as backend storage

There is only 3 steps before finally using remote ChubaoFS cluster as file system

1. Create StorageClass
2. Create PVC (Persistent Volume Claim)
3. Reference PVC in a Pod

### Create StorageClass

An example storage class yaml file is shown below.

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: cfs-sc
provisioner: csi.chubaofs.com
reclaimPolicy: Delete
```

Creating command.

```
$ kubectl create -f ~/storageclass-chubaofs.yaml
```

The field `provisioner` indicates name of the CSI driver, which is `csi.chubaofs.com` in this example. It is connected to the `drivername` in `deploy/csi-controller-deployment.yaml` and `deploy/csi-node-daemonset.yaml`. So StorageClass knows which driver should be used to manipulate the backend storage cluster.

| Name       | Madotory | Description|
| :--------- | :------: | ---------: |
| MasterAddr | Y | Master address of a specific on-premise ChubaoFS cluster |
| consulAddr | N | |


### Create PVC

An example pvc yaml file is shown below.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: cfs-pvc
spec:
  accessModes:
  - ReadWriteMany
  volumeMode: Filesystem
  resources:
    requests:
      storage: 2Gi
  storageClassName: cfs-sc
```

```
$ kubectl create -f ~/cfs-pvc.yaml
```

The field `storageClassName` refers to the StorageClass we already created.

### Use PVC in a Pod

The example `nginx-deployment.yaml` looks like below.

```yaml
...
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chubaofs-csi-demo
  labels:
    app: chubaofs-csi-demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: chubaofs-csi-demo-pod
  template:
    metadata:
      labels:
        app: chubaofs-csi-demo-pod
    spec:
      containers:
        - name: chubaofs-csi-demo
          image: nginx
          imagePullPolicy: "IfNotPresent"
          ports:
            - containerPort: 80
              name: "http-server"
          volumeMounts:
            - mountPath: "/usr/share/nginx/html"
              name: mypvc
      volumes:
        - name: mypvc
          persistentVolumeClaim:
            claimName: cfs-pvc
...
```

The field `claimName` refers to the PVC created before.
```
$ kubectl create -f ~/pv-pod.yaml
```
