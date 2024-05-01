#!/bin/bash
LAB_COMPOSE_FILE=deploy/docker-compose.yml
WEB_COMPOSE_FILE=deploy/web-service/docker-compose.yml

COMMAND=$0

usage () {
	echo "Usage $COMMAND [-h] [-f <path-to-compose-file>] [web]*."
	echo "(*default all containers)"
	exit
}

if [ "${1}" = "-h" ]
then
	usage
fi

if [ "${1}" = "-f" ]
then
	strip
	LAB_COMPOSE_FILE=${1}
	strip
fi

CONTAINERS=$*
if [ "${CONTAINERS}." = "." ]
then
	CONTAINERS="web"
fi

if [ ! -f ${LAB_COMPOSE_FILE} ]
then
	echo "Docker compose file not found."
	usage
fi

. $(dirname $0)/common.sh

for CNTR in $CONTAINERS
do
	if [ "${CNTR}" = "web" -a -f ${WEB_COMPOSE_FILE} ]
	then
		SVCHOME=$(dirname ${WEB_COMPOSE_FILE})
		pushd ${SVCHOME} > /dev/null
		docker-compose -f $(basename ${WEB_COMPOSE_FILE}) down
		docker rm -f df-labs-web-service > /dev/null 2>&1
		docker rmi df-labs/web-service > /dev/null 2>&1
		popd > /dev/null
	fi
done

docker-compose -f ${LAB_COMPOSE_FILE} down
