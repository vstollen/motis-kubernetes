# MOTIS Kubernetes

A project aiming to run [MOTIS](https://motis-project.de/) in Kubernetes.

## Parts of this Repository

### Operator
The operator is the main part of the project. It contains a Kubernetes operator which
manages the preprocessing and startup of MOTIS in Kubernetes.

### Init Container
The init container is used by the operator to setup the volumes and
download map data and train schedules needed by MOTIS.

### Kubernetes Objects
The "Kubernetes Objects" folder contains YAML files describing MOTIS-related
Kubernetes objects, mainly used for testing.
