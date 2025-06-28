#!/bin/bash

cd ./go-am-discord-rpc || {
  echo "Directory ./go-am-discord-rpc does not exist. Exiting."
  exit 1
}

echo "Building Go dependancies"
make "$@"
