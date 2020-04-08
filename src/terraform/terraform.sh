#!/bin/bash

OPERATION=$1
BACKEND=$2
FORCE=$3
OPTS="-var-file tfvars"

if [[ "$FORCE" = "true" ]]; then
  OPTS="-auto-approve $OPTS"
fi

if [[ "$OPERATION" = "show" ]]; then
  OPTS=""
fi

cd /app/$BACKEND
cp /data/* .
terraform init
terraform $OPERATION $OPTS
