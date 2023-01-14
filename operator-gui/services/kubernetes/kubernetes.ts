import {
  ApiextensionsV1Api,
  CoreV1Api,
  CustomObjectsApi,
  KubeConfig,
} from "@kubernetes/client-node";

const kubeConfig = new KubeConfig();
kubeConfig.loadFromDefault();

export const coreV1Api = kubeConfig.makeApiClient(CoreV1Api);
export const apiextensionsV1Api = kubeConfig.makeApiClient(ApiextensionsV1Api);
export const customObjectsApi = kubeConfig.makeApiClient(CustomObjectsApi);
