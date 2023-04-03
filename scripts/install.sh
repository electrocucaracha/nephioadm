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
if [[ ${DEBUG:-false} == "true" ]]; then
    set -o xtrace
    export PKG_DEBUG=true
fi

# shellcheck source=./scripts/_utils.sh
source _utils.sh

export PKG_KREW_PLUGINS_LIST=" "
export PKG_COMMANDS_LIST="docker,kubectl,git"
export PKG="go-lang"
export KPT_VERSION=1.0.0-beta.29

# NOTE: Shorten link -> https://github.com/electrocucaracha/pkg-mgr_scripts
curl -fsSL http://bit.ly/install_pkg | bash

if ! command -v kpt >/dev/null; then
    curl -s "https://i.jpillora.com/GoogleContainerTools/kpt@v${KPT_VERSION}!" | bash
    kpt completion bash | sudo tee /etc/bash_completion.d/kpt >/dev/null
fi

if ! command -v multicluster >/dev/null; then
    # shellcheck disable=SC1091
    [ -f /etc/profile.d/path.sh ] && source /etc/profile.d/path.sh
    GOBIN=/usr/local/bin/ sudo -E "$(command -v go)" install github.com/electrocucaracha/multicluster/cmd/multicluster@latest
fi
