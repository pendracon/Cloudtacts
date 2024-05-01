#!/usr/bin/env bash
. $(dirname $0)/../common.sh
MYSQL_ROOT_PASS=${MYSQL_ROOT_PASS:-${DEFPASS}}

docker run -d -p 33060:3306 --name mysql \
-v ${HOMEDIR}/mysql:/home/mysql \
-e "MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASS}" \
mysql:latest
