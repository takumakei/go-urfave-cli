#!/usr/bin/env bash
set -eu -o pipefail

main() {
  local i d
  while IFS= read -r -d $'\0' i; do
    d="$(dirname "$i")"
    echo "===> [$d]"
    pushd "$d"
    go get -u && go mod tidy
    popd
  done < <(find . -name go.mod -print0)
}

main
