#!/bin/sh

docker run --rm -it \
  --env-file "$(pwd)/.env" \
  -p 5556:5556 \
  -p 5557:5557 \
  -p 5558:5558 \
  --user 1000:1000 \
  --name dex \
  -v "$(pwd)/config.yaml:/config/config.yaml:ro" \
  ghcr.io/dexidp/dex:v2.37.0-distroless \
  dex serve /config/config.yaml
