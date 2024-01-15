#!/usr/bin/env bash

# How to run manually:
# docker build --pull --rm -f "contrib/devtools/Dockerfile" -t cosmossdk-proto:latest "contrib/devtools"
# docker run --rm -v $(pwd):/workspace --workdir /workspace cosmossdk-proto sh ./scripts/protocgen.sh

echo "Formatting protobuf files"
find ./ -name "*.proto" -exec clang-format -i {} \;

set -e

echo "Generating gogo proto code"
cd ./proto
proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  proto_files=$(find "${dir}" -maxdepth 1 -name '*.proto')
  for file in $proto_files; do
    if grep -q "option go_package" "$file"; then
      buf generate --template buf.gen.gogo.yaml "$file"
    fi
  done
done

cd ..

cp -r github.com/settlus/chain/evmos/* ./evmos
cp -r github.com/settlus/chain/x/* ./x
rm -rf github.com
