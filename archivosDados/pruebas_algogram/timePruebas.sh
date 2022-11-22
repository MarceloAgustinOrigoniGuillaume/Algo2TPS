#!/usr/bin/env bash

set -eu

PROGRAMA="$1"

echo "running: ./pruebas.sh $1"
START_TIME=$(date +%s)

./pruebas.sh $1

END_TIME=$(date +%s)

echo "It took $(($END_TIME-$START_TIME)) seconds"
