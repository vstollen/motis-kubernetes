import k8s from "@kubernetes/client-node";

const kubeConfig = new k8s.KubeConfig();
kubeConfig.loadFromDefault();

export const coreV1Api = kubeConfig.makeApiClient(k8s.CoreV1Api);
export const apiextensionsV1Api = kubeConfig.makeApiClient(k8s.ApiextensionsV1Api);
export const customObjectsApi = kubeConfig.makeApiClient(k8s.CustomObjectsApi);
