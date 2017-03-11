#!/bin/sh

set -e

export REPO_DIR="$(cd "$(dirname "${0}")/../.." && pwd)"
export TMPL_DIR="etc/tmpl"
export GO_TMPL_FILES="google/protobuf/protobuf.gen.go.tmpl,google/type/type.gen.go.tmpl,pb/time/time.gen.go.tmpl"

go run "${REPO_DIR}/etc/cmd/gen-blank/main.go"
