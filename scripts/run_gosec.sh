#!/bin/bash
# This script runs the gosec scanner locally

if ! command -v gosec 2> /dev/null
then
  echo "error gosec must be installed with this command: make gosec_install" && exit 1
fi

gosec -no-fail -fmt=sarif -out=gosec.sarif -exclude-dir test  -exclude-dir generator  ./...
