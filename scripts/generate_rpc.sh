#!/bin/bash

set -euxo pipefail

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

protoc -I $SCRIPTDIR/../rpc $SCRIPTDIR/../rpc/types/*.proto --go_out=plugins=grpc:$(go env GOPATH)/src
protoc -I $SCRIPTDIR/../rpc $SCRIPTDIR/../rpc/mempool/*.proto --go_out=plugins=grpc:$(go env GOPATH)/src
protoc -I $SCRIPTDIR/../rpc $SCRIPTDIR/../rpc/admission_control/*.proto --go_out=plugins=grpc:$(go env GOPATH)/src
