#!/bin/bash
# SPDX-license-identifier: Apache-2.0
##############################################################################
# Copyright (c) 2023
# All rights reserved. This program and the accompanying materials
# are made available under the terms of the Apache License, Version 2.0
# which accompanies this distribution, and is available at
# http://www.apache.org/licenses/LICENSE-2.0
##############################################################################

set -o pipefail
set -o errexit
set -o nounset
[[ ${DEBUG:-false} != "true" ]] || set -o xtrace

# shellcheck source=./scripts/_utils.sh
source _utils.sh

# shellcheck source=./scripts/_common.sh
source _common.sh

trap get_status ERR

base_path=/tmp/nephio_pkgs
sudo mkdir -p "$base_path"/{mgmt,edge}

# Multi-cluster configuration
if ! sudo docker ps --format "{{.Image}}" | grep -q "kindest/node"; then
    sudo multicluster create --config ./config.yml --name nephio
    mkdir -p "$HOME/.kube"
    sudo cp /root/.kube/config "$HOME/.kube/config"
    sudo chown -R "$USER": "$HOME/.kube"
    chmod 600 "$HOME/.kube/config"
fi

# shellcheck disable=SC1091
[ -f /etc/profile.d/path.sh ] && source /etc/profile.d/path.sh
pushd "$(git rev-parse --show-toplevel)" >/dev/null

# Bootstrap management cluster
sudo kubectl config use-context kind-nephio
sudo "$(command -v go)" run ./... init --base-path "$base_path/mgmt" --debug \
    --git-service "http://$(ip route get 8.8.8.8 | grep '^8.' | awk '{ print $7 }'):3000/nephio-playground" \
    --backend-base-url "http://localhost:7007" \
    --webui-cluster-type NodePort

# Join additional clusters
sudo kubectl config use-context kind-regional
sudo "$(command -v go)" run ./... join --base-path "$base_path/edge" --debug \
    --git-service "http://$(ip route get 8.8.8.8 | grep '^8.' | awk '{ print $7 }'):3000/nephio-playground"
popd >/dev/null
