#!/usr/bin/env bash
. $(dirname $0)/../common.sh
FIRESTORE_PROJECT_ID=${FIRESTORE_PROJECT_ID:-${DEFPROJECT}}

docker run -d --name firebase-suite \
-v ${HOMEDIR}:/home/node \
--port 4000:4000 --port 5000:5000 --port 5001:5001 \
--port 8080:8080 --port 8085:8085 --port 9000:9000 \
--port 9005:9005 --port 9099:9099 --port 9199:9199 \
andreysenov/firebase-tools firebase emulators:start
