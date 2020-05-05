#!/bin/bash
# Reference: https://coreos.com/quay-enterprise/docs/latest/mysql-container.html
set -e

MYSQL_USER="LOTRuser"
MYSQL_DATABASE="LOTRdata"
MYSQL_CONTAINER_NAME="mysqlshire"
LOCAL_DB_DIR=~/mydb/mysql-datadir
HOST_PORT=13306

# for better passwords, use
#  $(uuidgen | sed "s/-//g") or $(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | sed 1q)
MYSQL_ROOT_PASSWORD=$(echo LOTRrootpass)
MYSQL_PASSWORD=$(echo LOTRpass)

mkdir -p ${LOCAL_DB_DIR}
echo "Starting the MySQL container as '${MYSQL_CONTAINER_NAME}'"

docker \
  run \
  --detach \
  --env MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
  --env MYSQL_USER=${MYSQL_USER} \
  --env MYSQL_PASSWORD=${MYSQL_PASSWORD} \
  --env MYSQL_DATABASE=${MYSQL_DATABASE} \
  --name ${MYSQL_CONTAINER_NAME} \
  --volume ${LOCAL_DB_DIR}:/var/lib/mysql \
  --publish ${HOST_PORT}:3306 \
  mysql;

echo "Database '${MYSQL_DATABASE}' running."
echo "  Username: ${MYSQL_USER}"
echo "  Password: ${MYSQL_PASSWORD}"
echo "Port ${HOST_PORT}"
echo "Persisting to local directory ${LOCAL_DB_DIR}"
