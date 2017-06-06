#!/bin/bash

# Terminate script if any simple command fails
set -e

# Check if environment given
if [ -z ${1+x} ]
then
  echo "Deployment destination is not set: indicate 'local', 'dev', 'prod' as an argument."
  exit 1
fi

# Check if GOPATH set
if [ -z ${GOPATH+x} ]
then
  echo "GOPATH needs to be set -- export GOPATH=_____"
  exit 1
fi

# Check if GOBIN set
if [ -z ${GOBIN+x} ]
then
  echo "GOBIN needs to be set -- export GOBIN=_____"
  exit 1
fi

# Set s3 target bucket values
if [ "$1" == "local" ]
then
  TARGETBUCKET="target-bucket-dev"
  TARGETOBJECT="goMicroserviceTemplate/local/goMicroserviceTemplate"
elif [ "$1" == "dev" ]
then
  TARGETBUCKET="target-bucket-dev"
  TARGETOBJECT="goMicroserviceTemplate/goMicroserviceTemplate"
elif [ "$1" == "prod" ]
then
  TARGETBUCKET="target-bucket-prod"
  TARGETOBJECT="goMicroserviceTemplate/goMicroserviceTemplate"
  export GIN_MODE=release
else
  echo "Deployment destination must be 'local', 'dev', or 'prod'"
  exit 1
fi

echo "### Running go install ###"
go install *.go

echo "### Compressing Binary ###"
gzip -kfn $GOBIN/goMicroserviceTemplate

echo "### Uploading Binary to $TARGETBUCKET/$TARGETOBJECT.gz ###"
# Check if the MD5 sum of the package is different than what we have
if test -e $GOBIN/goMicroserviceTemplate.gz && aws s3api head-object --bucket $TARGETBUCKET --key $TARGETOBJECT.gz | grep `md5sum $GOBIN/goMicroserviceTemplate.gz | awk '{print $1;}'` > /dev/null; then
  # Found, do nothing
  echo "Latest build already deployed, quitting"
else
  echo "Uploading to s3 for $1 server(s)"
  aws s3api put-object --bucket $TARGETBUCKET --key $TARGETOBJECT.gz --body $GOBIN/goMicroserviceTemplate.gz
fi

echo "### Build complete ###"
