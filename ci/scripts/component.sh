#!/bin/bash -eux

pushd dp-integrity-checker
  make test-component
popd
