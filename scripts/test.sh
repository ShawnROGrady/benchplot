#! /usr/bin/env bash

set -eu -o pipefail

# test go code
go test -race -cover ./...

# check if README up to date
./scripts/generate_readme.sh "TMP_README.md"
diff README.md TMP_README.md || (>&2 echo "README out of date"; rm "TMP_README.md"; exit 1)
rm "TMP_README.md"
