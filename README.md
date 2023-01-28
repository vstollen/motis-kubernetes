# MOTIS Kubernetes

A project helping to run [MOTIS](https://motis-project.de/) in Kubernetes.

## Parts of this Repository
### Operator
The operator is the main part of the project. It contains a Kubernetes operator which
manages the preprocessing and startup of MOTIS in Kubernetes.

### Operator GUI
The operator GUI is a web application which allows users to create and manage MOTIS instances using a simple GUI.
Therefore, users do not need Kubernetes knowledge to manage their MOTIS instances.

### Init Container
The init container is used by the operator to setup the volumes and
download map data and train schedules needed by MOTIS.

### Kubernetes Objects
The "Kubernetes Objects" folder contains MOTIS-related Kubernetes objects, mainly used for testing.

## Setup
### Prerequisites
You’ll need a Kubernetes cluster to run the operator against. You can use [Minikube](https://minikube.sigs.k8s.io/docs/) to get a local cluster for testing, or run against a remote cluster.
**Note:** The controller and the Operator GUI will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Operator
Visit the [Operator `README.md`](motis-operator/README.md) for instructions on how to set up the Operator.

### Operator GUI
Visit the [Operator GUI `README.md`](operator-gui/README.md) for instructions on how to set up the Operator GUI.
While the Operator GUI requires the Operator to run, you do not need to run the Operator GUI to use the Operator.