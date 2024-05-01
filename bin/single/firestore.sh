#!/usr/bin/env bash
. $(dirname $0)/../common.sh
FIRESTORE_PROJECT_ID=${FIRESTORE_PROJECT_ID:-${DEFPROJECT}}

docker run -d -p 8200:8200 --name firebase-emulator \
--env "FIRESTORE_PROJECT_ID=${FIRESTORE_PROJECT_ID}" \
--env "PORT=8200" \
mtlynch/firestore-emulator-docker
