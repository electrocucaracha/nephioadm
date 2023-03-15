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
set -o xtrace

# shellcheck source=./scripts/_utils.sh
source _utils.sh

export DEBUG=true

for step in install configure; do
    info "Running $step process"
    bash "./$step.sh"
    [[ ${ENABLE_FUNC_TEST:-false} != "true" ]] || bash "./${step}_test.sh"
done
