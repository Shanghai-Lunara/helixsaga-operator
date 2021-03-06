#!/usr/bin/env bash

export GOPATH=`go env | grep -i gopath | awk '{split($0,a,"\""); print a[2]}'`

# 代码生成的工作目录，也就是我们的项目路径
ROOT_PACKAGE="github.com/Shanghai-Lunara/helixsaga-operator"
# API Group
CUSTOM_RESOURCE_NAME="helixsaga"
# API Version
CUSTOM_RESOURCE_VERSION="v1"

GENS="$1"

if [ "${GENS}" = "crd" ] || grep -qw "crd" <<<"${GENS}"; then
  cp ${GOPATH}/bin/go-to-protobuf-crd ${GOPATH}/bin/go-to-protobuf
  Packages="$ROOT_PACKAGE/pkg/apis/$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION"
  "${GOPATH}/bin/go-to-protobuf" \
     --packages "${Packages}" \
     --clean=false \
     --only-idl=false \
     --keep-gogoproto=false \
     --verify-only=false \
     --proto-import ${GOPATH}/src/k8s.io/api/core/v1
fi

# 执行代码自动生成，其中pkg/client是生成目标目录，pkg/apis是类型定义目录
${GOPATH}/src/k8s.io/code-generator/generate-groups.sh all "$ROOT_PACKAGE/pkg/generated/$CUSTOM_RESOURCE_NAME" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
