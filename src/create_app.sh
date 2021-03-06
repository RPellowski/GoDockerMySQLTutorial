#!/bin/bash
. my-env

docker build -t hobbit -f Dockerfile.app --network host .

docker run \
    --detach \
    --publish 8082:8080 \
    --name frodo \
    --env-file my-env \
    --network ${APP_NETWORK} \
    hobbit
