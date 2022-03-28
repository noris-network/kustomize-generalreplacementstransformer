Use a kustomize [patchStrategicMerge](https://github.com/kubernetes-sigs/kustomize/blob/master/docs/glossary.md#patchstrategicmerge) to apply the following patches to patch the ArgoCD `install.yaml`.

### generalReplacementsTransformer.yaml

    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: argocd-repo-server
    spec:
      template:
        spec:
          volumes:
            - name: custom-tools
              emptyDir: {}
          initContainers:
            - name: install-generalreplacementstransformer
              image: alpine:3.15
              command: ["/bin/sh", "-c"]
              args:
                - |
                  echo "Installing noris.net/v1alpha1/GeneralReplacementsTransformer..."
                  set -ex
                  wget -O /custom-tools/GeneralReplacementsTransformer https://github.com/noris-network/kustomize-generalreplacementstransformer/releases/download/v${VERSION}/GeneralReplacementsTransformer_${VERSION}_${OS}_${ARCH}
                  chmod -v +x /custom-tools/GeneralReplacementsTransformer
                  set +x
                  echo "Done."
              volumeMounts:
                - mountPath: /custom-tools
                  name: custom-tools
              env:
                - name: VERSION
                  value: 0.10.1
                - name: OS
                  value: linux
                - name: ARCH
                  value: amd64
          containers:
            - name: argocd-repo-server
              volumeMounts:
                - mountPath: /.config/kustomize/plugin/noris.net/v1alpha1/generalreplacementstransformer/GeneralReplacementsTransformer
                  name: custom-tools
                  subPath: GeneralReplacementsTransformer
              env:
                - name: XDG_CONFIG_HOME
                  value: /.config

### enablePlugins.yaml

    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: argocd-cm
    data:
      kustomize.buildOptions: --enable_alpha_plugins
