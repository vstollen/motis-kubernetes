# MOTIS Operator GUI
The Operator GUI is a webapp that allows to create MOTIS instances in a Kubernetes Cluster where the MOTIS operator is
deployed.

## Setup
### Prerequisites
You’ll need a Kubernetes cluster to run against. You can use [Minikube](https://minikube.sigs.k8s.io/docs/) to get a local cluster for testing, or run against a remote cluster.
**Note:** The operator GUI will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

Also, you need do deploy the MOTIS Operator to that Kubernetes cluster. For instructions on how to do that,
refer to the [MOTIS Operator `README.md`](../motis-operator/README.md).

### Running the Operator GUI
To run the operator GUI, you need to build and start the Operator GUI server:
```bash
npm run build
npm run start
```

## Development

First, run the development server:

```bash
npm run dev
# or
yarn dev
```

## Learn More

The Operator GUI uses React.js and Next.js. To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js/) - your feedback and contributions are welcome!
