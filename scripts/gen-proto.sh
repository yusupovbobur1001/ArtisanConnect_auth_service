#!/bin/bash
CURRENT_DIR=$1
rm -rf ${CURRENT_DIR}/genproto
for x in $(find ${CURRENT_DIR}/ArtisanConnect_protos1 -type f -name '*.proto'); do
  protoc -I=$(dirname ${x}) -I=${CURRENT_DIR}/ArtisanConnect_protos1 --go_out=${CURRENT_DIR} \
   --go-grpc_out=${CURRENT_DIR} ${x}
done
