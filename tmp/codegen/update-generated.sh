#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
gitlab.com/pickledrick/vpc-peering-operator/pkg/generated \
gitlab.com/pickledrick/vpc-peering-operator/pkg/apis \
r4:v1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
