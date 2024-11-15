# Yggdrasil

Mono repo setup for my Raspberry Pi 4 with k3s setup. Experimenting with a,
mostly self hosted, platform setup.

## Requirements

- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [taskfile](https://taskfile.dev/installation/)
- [helm](https://helm.sh/docs/intro/quickstart/)
- [helmfile](https://helmfile.readthedocs.io/en/latest/#installation)
- [gcloud](https://cloud.google.com/sdk/docs/install)

You will also need a Google Cloud project and appropriate access to make this
setup work out of the box. Should work with other cloud providers with a few
changes here and there.

## Taskfiles

I try to keep all needed scripts/tasks/actions in the various taskfiles and
make them as reusable as possible. Eg. since there are potentially multiple
applications in this repo, the go and docker taskfiles use wildcard to specify
the application. So, to build the experimental application docker, you would
run `task docker:build-experimental` and then push it with `task
docker:push-experimental`.

## Raspberry Pi 4 Setup

Start with a fresh install of Raspberry Pi OS Lite.

```bash
sudo apt-get update && sudo apt-get upgrade -y
```

## k3s

### Installation

Append `cgroup_memory=1 cgroup_enable=memory` to `/boot/firmware/cmdline.txt`
and do a reboot.

Run `curl -sfL https://get.k3s.io | sh -`. Which is used for installing and
upgrading.

Follow the instructions in [Cluster
Access](https://docs.k3s.io/cluster-access). I'm not a fan of this way of
authenticating, but it should be relatively safe as long as proper security
measures are in place. It works as long as there is only one user.

Follow the [hardening guide](https://docs.k3s.io/security/hardening-guide).

You will need to add a rule that allows Prometheus to scrape kube-dns as well:

```yaml

apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: kube-dns-allow-prometheus
  namespace: kube-system
spec:
  podSelector:
    matchLabels:
      k8s-app: kube-dns
  policyTypes:
  - Ingress
  ingress:
    - from:
      - namespaceSelector:
          matchLabels:
            kubernetes.io/metadata.name: monitoring
      - podSelector:
          matchLabels:
            app: prometheus
      ports:
        - protocol: TCP
          port: 9153
```

### Access to Private Repository

Info on how to set up private registries can be found
[here](https://docs.k3s.io/installation/private-registry). It will probably
vary a bit depending on what repository one uses. In my case I'm using Google
Artifact Registry.

First step is to create a service account with `roles/artifactregistry.reader`.
Then create a key. Use this key in the password field below, but put it into
one line. Also remember to change `europe-north1-docker.pkg.dev` to match your
setup.

You can run `create-docker-read-service-account` to set up the service account,
attach the role and create a key. Replace `<YOUR_SERVICE_ACCOUNT_KEY>` with the
contents of the key file. Remove newlines before putting it in.

```yaml
configs:
  europe-north1-docker.pkg.dev:
    auth:
      username: _json_key
      password: '<YOUR_SERVICE_ACCOUNT_KEY>'
```

This needs to be done on all nodes if you have more.

### Deploying

Install [helm](https://helm.sh/docs/intro/install/) on your local machine.

Install [helmfile](https://helmfile.readthedocs.io/en/latest/#installation) on
your local machine. The reason for using helmfile is to simplify deploying
everything in one go. Which gets difficult when deploying to multiple
namespaces.

Run `task helmfile:apply` to deploy everything.
