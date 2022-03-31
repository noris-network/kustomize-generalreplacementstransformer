# kustomize-generalreplacementstransformer

[![Go Report Card](https://goreportcard.com/badge/github.com/noris-network/kustomize-generalreplacementstransformer)](https://goreportcard.com/report/github.com/noris-network/kustomize-generalreplacementstransformer)
[![Latest Release](https://img.shields.io/github/v/release/noris-network/kustomize-generalreplacementstransformer?sort=semver)](https://github.com/noris-network/kustomize-generalreplacementstransformer/releases/latest)
[![License](https://img.shields.io/github/license/noris-network/kustomize-generalreplacementstransformer)](https://github.com/noris-network/kustomize-generalreplacementstransformer/blob/main/LICENSE)

## What is this for?
[Kustomize](https://github.com/kubernetes-sigs/kustomize) is a great tool
for deploying Applications following GitOps. But Sometimes you need to
change "things" that are not addressable with the build in
[replacements](https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/replacements/).
That's where GeneralReplacementsTransformer comes into play. It's a kustomize
plugin that allows you to select values in a similar way than the build in
replacements, but uses golang template expressions wherever you need to insert
values. This is very powerful, but should definitely be used with care.

## Installation

The `GeneralReplacementsTransformer` binary can be downloaded from the
[GitHub releases page](https://github.com/noris-network/kustomize-generalreplacementstransformer/releases).
In order to be called by [kustomize](https://github.com/kubernetes-sigs/kustomize),
it has to be installed to `$XDG_CONFIG_HOME/kustomize/plugin/noris.net/v1alpha1/generalreplacementstransformer`.
(`$XDG_CONFIG_HOME` points by default to `$HOME/.config` on Linux and OS X, and `%LOCALAPPDATA%` on Windows.)

Install version 0.11.1 on Linux:

    VERSION=0.11.1 OS=linux ARCH=amd64
    INSTALL_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/kustomize/plugin/noris.net/v1alpha1/generalreplacementstransformer"
    curl -Lo GeneralReplacementsTransformer https://github.com/noris-network/kustomize-generalreplacementstransformer/releases/download/v${VERSION}/GeneralReplacementsTransformer_${VERSION}_${OS}_${ARCH}
    chmod +x GeneralReplacementsTransformer
    mkdir -p $INSTALL_DIR
    mv GeneralReplacementsTransformer $INSTALL_DIR

## Usage

Let's say you need a password in more than one place, but some locations are not
addressable by build in replacements, and you only want to define it once...

Create a kustomization.yaml file:

    cat <<. >kustomization.yaml
    apiVersion: kustomize.config.k8s.io/v1beta1
    kind: Kustomization
    namespace: demo
    secretGenerator:
      - name: mongodb-auth
        literals:
          - mongodb-root-password=secret123
      - name: mongodb-env
        literals:
          - MONGO_URL=mongodb://demo:{{.password}}@mongodb/demo
    transformers:
      - transformer.yaml
    .

    cat <<. >transformer.yaml
    apiVersion: noris.net/v1alpha1
    kind: GeneralReplacementsTransformer
    metadata:
      name: example
    selectValues:
      - name: password
        resource:
          kind: Secret
          name: mongodb-auth
          fieldPath: data.mongodb-root-password
    replacements:
      - resource:
          kind: Secret
          name: mongodb-env
        type: template
    .

    kustomize build --enable-alpha-plugins

It is of cause not recommended to put your secret data unencrypted into any files,
you could e.g. use [SopsSecretGenerator](https://github.com/goabout/kustomize-sopssecretgenerator)
to protect them. GeneralReplacementsTransformer will still work.

## Selecting Values

The `resource`-selector in `selectValues` supports `kind`, `name` and `fieldPath`.

## Inserting Values

The `resource`-selector in `replacements` supports `kind` and `name`, which might
be empty to select multiple resources.

All string values in yaml content can contain golang template expressions, e.g.:

    key: "{{.value}}"

[Slim-sprig](https://go-task.github.io/slim-sprig/) function are also available:

    key: "deployed at {{ now | date "2006-01-02 }}"

Right now just `type: template` is supported, this might change some time, but there
are no plans so far.

## Using GeneralReplacementsTransformer with ArgoCD

GeneralReplacementsTransformer can be added to ArgoCD by [patching](doc/argocd.md)
an initContainer into the ArgoCD provided `install.yaml`.
