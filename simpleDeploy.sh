# This is an altnerate, simpler deployment method for when the service will be run on a single vm

# Terminate script if any simple command fails
set -e

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

# Check if environment given
if [ -z ${1+x} ]
then
  echo "Deployment destination is not set: indicate 'local', 'dev', or 'prod' as an argument."
  exit 1
fi

if [ "$1" == "local" ]
then
  export GIN_MODE=debug
elif [ "$1" == "dev" ]
then
  export GIN_MODE=debug
  TARGETVM="vm.domain.or.ip"
  KEYLOCATION="/path/to/dev/ssh/key"
elif [ "$1" == "prod" ]
then
  TARGETVM="vm.domain.or.ip"
  KEYLOCATION="/path/to/prod/ssh/key"
  export GIN_MODE=release
else
  echo "Deployment destination must be 'local', 'dev', or 'prod'"
  exit 1
fi

echo "### Running go install ###"
go install *.go

if [ "$1" == "local" ]
then
  sudo cp goMicroserviceTemplate.conf /etc/init/
  sudo service goMicroserviceTemplate stop || true
  sudo cp $GOBIN/goMicroserviceTemplate /usr/local/bin/
  sudo service goMicroserviceTemplate start
  exit 0
fi

gzip -f $GOBIN/goMicroserviceTemplate

scp -i $KEYLOCATION goMicroserviceTemplate.conf ubuntu@$TARGETVM:/home/ubuntu/
scp -i $KEYLOCATION $GOBIN/goMicroserviceTemplate.gz ubuntu@$TARGETVM:/home/ubuntu/

ssh ubuntu@$TARGETVM -i $KEYLOCATION 'sudo mv goMicroserviceTemplate.conf /etc/init/'
ssh ubuntu@$TARGETVM -i $KEYLOCATION 'gzip -d goMicroserviceTemplate.gz; sudo mv goMicroserviceTemplate /usr/local/bin/; sudo service goMicroserviceTemplate restart'
