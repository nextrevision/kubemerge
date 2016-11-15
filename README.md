# kubemerge

Merge two Kubernetes resources together.

kubemerge was written to help inject partials into existing Kubernetes resources without the need for templating. For example, if I wanted to enforce that imagePullSecrets was set for a deployment:

```
# fixtures/source.yml
---
spec:
  template:
    spec:
      imagePullSecrets:
        - name: registry
```

```
# fixtures/deployment.yml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:latest
          ports:
            - containerPort: 80
```

Then I can merge them with kubemerge:

```
$ kubemerge -yaml fixtures/source.yml fixtures/deployment.yml
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: nginx
  name: nginx
spec:
  replicas: 3
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx:latest
        name: nginx
        ports:
        - containerPort: 80
        resources: {}
      imagePullSecrets:              # note the merged values
      - name: registry               # from the source file
status: {}
```

## Usage

```
Usage: kubemerge [options] source target

Examples:

  kubemerge policy.yaml deployment.yaml
  kubemerge -yaml policy.yaml deployment.yaml
	cat deployment.yaml | kubemerge policy.yaml

Options:
  -yaml
    	Output as YAML
```

## Supported Resources

* Deployment
* DaemonSet
* ReplicaSet
* ReplicationController
