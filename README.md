# Nephio installer
<!-- markdown-link-check-disable-next-line -->
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub Super-Linter](https://github.com/electrocucaracha/nephioadm/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)
[![Ruby Style Guide](https://img.shields.io/badge/code_style-rubocop-brightgreen.svg)](https://github.com/rubocop/rubocop)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrocucaracha/nephioadm)](https://goreportcard.com/report/github.com/electrocucaracha/nephioadm)
[![GoDoc](https://godoc.org/github.com/electrocucaracha/nephioadm?status.svg)](https://godoc.org/github.com/electrocucaracha/nephioadm)
![visitors](https://visitor-badge.glitch.me/badge?page_id=electrocucaracha.nephioadm)

This tool provisions [Nephio][1] components on target clusters

```bash
go install github.com/electrocucaracha/nephioadm/...
nephioadm init --base-path "/opt/nephio/mgmt" --git-service "http:/gitea-server:3000/nephio-playground" 
```

## Provisioning process

The [Vagrant tool][2] can be used for provisioning an Ubuntu Focal
Virtual Machine. It's highly recommended to use the  *setup.sh* script
of the [bootstrap-vagrant project][3] for installing Vagrant
dependencies and plugins required for this project. That script
supports two Virtualization providers (Libvirt and VirtualBox) which
are determine by the **PROVIDER** environment variable.

```bash
curl -fsSL http://bit.ly/initVagrant | PROVIDER=libvirt bash
```

Once Vagrant is installed, it's possible to provision a Virtual
Machine using the following instructions:

```bash
vagrant up
```

[1]: https://nephio.org/
[2]: https://www.vagrantup.com/
[3]: https://github.com/electrocucaracha/bootstrap-vagrant
