# Nephio installer
<!-- markdown-link-check-disable-next-line -->
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub Super-Linter](https://github.com/electrocucaracha/nephioadm/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)
[![Ruby Style Guide](https://img.shields.io/badge/code_style-rubocop-brightgreen.svg)](https://github.com/rubocop/rubocop)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrocucaracha/nephioadm)](https://goreportcard.com/report/github.com/electrocucaracha/nephioadm)
[![GoDoc](https://godoc.org/github.com/electrocucaracha/nephioadm?status.svg)](https://godoc.org/github.com/electrocucaracha/nephioadm)
![visitors](https://visitor-badge.glitch.me/badge?page_id=electrocucaracha.nephioadm)

This tool provisions [Nephio][1] components on target clusters

## Installation

```bash
go install github.com/electrocucaracha/nephioadm/cmd/nephioadm@latest
```

## Usage

For management components (system, webui and configsync packages):

```bash
nephioadm init \
    --base-path "/opt/nephio/mgmt" \
    --git-service "http:/gitea-server:3000/nephio-playground" \
    --backend-base-url "http://localhost:7007" \
    --webui-cluster-type NodePort
```

For workload components (configsync package):

```bash
nephioadm join \
    --base-path "/opt/nephio/mgmt" \
    --git-service "http:/gitea-server:3000/nephio-playground" 
```

## Provisioning process

This process uses two main components:

* Nephio packages ([official repository][1] by default). The `--nephio-repo`
argument allows the consumption of other sources. This can be useful during the
Nephio development and testing processes.
* Target clusters. Currently, this tool installs Nephio components on the
current pointing Kubernetes cluster. This cluster must be reachable from the
tool and requires the installation of [kpt CLI][2].

```text
          +-------------------------------------------------------+
          |   https://github.com/nephio-project/nephio-packages   |
          |    +----------+    +--------------+    +---------+    |
          |    |  system  |    |  configsync  |    |  webui  |    |
          |    +----------+    +--------------+    +---------+    |
          +-------------------------------------------------------+

                               +-----------+
                               | nephioadm |
                               +-----------+

+---------------------------------+     +---------------------------------+
| mgmt (k8s)                      |     | workload (k8s)                  |
| +-----------------------------+ |     | +-----------------------------+ |
| | mgmt-control-plane          | |     | | workload-control-plane      | |
| | podSubnet: 10.196.0.0/16    | |     | | podSubnet: 10.197.0.0/16    | |
| | serviceSubnet: 10.96.0.0/16 | |     | | serviceSubnet: 10.97.0.0/16 | |
| +-----------------------------+ |     | +-----------------------------+ |
+---------------------------------+     +---------------------------------+
```

The `--git-service` argument specifies the URL of the repository used by
[ConfigSync][3] and [Porch][4] components.


[1]: https://github.com/nephio-project/nephio-packages.git
[2]: https://kpt.dev/installation/kpt-cli
[3]: https://cloud.google.com/anthos-config-management/docs/config-sync-overview
[4]: https://kpt.dev/book/08-package-orchestration/
