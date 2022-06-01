#!/bin/bash

# this file is used for testing, set the DOMAIN below to your own name: DO NOT set it to beta and especially DO NOT
# set it to prod. We will use the Jenkins pipeline to deploy to these stages

# exit when any command fails
set -e

domain='john'

# delete build directory if exists
[ -d "./lib/build/" ] && rm -rf ./lib/build/

# build the go project
cwd=$(pwd)
cd ./lib/src/impruviservice/
GOARCH=amd64 GOOS=linux go build -o $cwd/lib/build/ImpruviService

# deploy to aws
cd $cwd
DOMAIN=$domain cdk synth
DOMAIN=$domain cdk deploy --require-approval never

# cleanup
rm -rf ./lib/build/
