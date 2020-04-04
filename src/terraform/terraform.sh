#!/bin/bash

OPERATION=$1
BACKEND=$2
FORCE=$3
OPTS=""

if [[ "$FORCE" = "true" ]]; then
  OPTS="-auto-approve $OPTS"
fi

cd /app/src/terraform/$BACKEND
terraform $OPERATION $OPTS
