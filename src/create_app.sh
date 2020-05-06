#!/bin/bash
docker build -t hobbit .
docker run --detach --publish 8082:8080 --name frodo --link shire:mysql hobbit
