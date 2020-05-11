#!/bin/bash
. my-env
docker build -t hobbit --network host .

docker run \
    --detach \
    --publish 8082:8080 \
    --name frodo \
    --env MYSQL_DATABASE=${MYSQL_DATABASE} \
    --env MYSQL_USER=${MYSQL_USER} \
    --env MYSQL_PASSWORD=${MYSQL_PASSWORD} \
    --env MYSQL_PORT=${MYSQL_PORT} \
    --env MYSQL_CONTAINER_NAME=${MYSQL_CONTAINER_NAME} \
    --network ${APP_NETWORK} \
    hobbit;
