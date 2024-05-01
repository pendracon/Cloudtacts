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

#-- Create application layout
#
for DIR in minio/data mysql/data
do
	if [ ! -d ${HOMEDIR}/${DIR} ]
	then
		mkdir -p ${HOMEDIR}/${DIR}
		chmod 777 ${HOMEDIR}/${DIR}
	fi
done

docker-compose -f ${LAB_COMPOSE_FILE} up -d

for CNTR in $CONTAINERS
do
	if [ "${CNTR}" = "web" -a -f ${WEB_COMPOSE_FILE} ]
	then
		SVCHOME=$(dirname ${WEB_COMPOSE_FILE})
		cp -uv bin/webserviceexe ${SVCHOME}/webserviceexe > /dev/null 2>&1
		pushd ${SVCHOME} > /dev/null
		docker rm -f df-labs-web-service > /dev/null 2>&1
		docker rmi df-labs/web-service > /dev/null 2>&1
		docker-compose -f $(basename ${WEB_COMPOSE_FILE}) up --build -d
		popd > /dev/null
	fi
done
