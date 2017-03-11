#!/bin/sh

set -e

export REPO_DIR="$(cd "$(dirname "${0}")/../.." && pwd)"
export TMPL_DIR="etc/tmpl"
export GO_TMPL_FILES="pb/money/money.gen.go.tmpl"
export PB_TMPL_FILES="pb/money/money.gen.proto.tmpl"
export CSV_FILE="etc/data/country-codes.csv"

go run "${REPO_DIR}/etc/cmd/gen-money/main.go"
