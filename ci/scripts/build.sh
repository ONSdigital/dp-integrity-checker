#!/bin/bash -eux

pushd dp-integrity-checker
  make build
  cp build/dp-integrity-checker Dockerfile.concourse ../build
popd
