#!/bin/bash
# Reference: https://coreos.com/quay-enterprise/docs/latest/mysql-container.html
set -e

. my-env

docker network inspect ${APP_NETWORK} >/dev/null 2>&1 || \
    docker network create ${APP_NETWORK}

mkdir -p ${LOCAL_DB_DIR}
echo "Starting the MySQL container as '${MYSQL_CONTAINER_NAME}'"

docker build -t shire -f Dockerfile.db .

docker \
  run \
  --detach \
  --env-file my-env \
  --name ${MYSQL_CONTAINER_NAME} \
  --volume ${LOCAL_DB_DIR}:/var/lib/mysql \
  --publish ${MYSQL_HOST_PORT}:3306 \
  --network ${APP_NETWORK} \
  shire;

echo "Database '${MYSQL_DATABASE}' running."
echo "  Username: ${MYSQL_USER}"
echo "  Password: ${MYSQL_PASSWORD}"
echo "Port ${MYSQL_PORT}"
echo "Persisting to local directory ${LOCAL_DB_DIR}"
