#!/bin/bash
pushd `dirname $0` > /dev/null
SCRIPT_PATH=`pwd -P`
popd > /dev/null
SCRIPT_FILE=`basename $0`


# msg color define
COLOR_RED='\033[0;31m'
COLOR_INFO='\033[0;36m'
COLOR_NC='\033[0m' # No Color

CONFIG_PATH=$SCRIPT_PATH"/meshservice"

if [ $# -eq 0 ] ; then
    echo "no provied eg server url, use default value: 127.0.0.1:2381"
    echo -e "${COLOR_INFO}usage: ${SCRIPT_FILE}  egservice-url ${COLOR_NC}"
fi

check_result ()
{
  if [ $1 !=  0 ]; then
      echo -e "${COLOR_RED}create easemesh tenant/service failed, errno: "$1${COLOR_NC}
      exit 1
  fi
}

EGCMD="egctl --server "$1

K8S_YAML_PATH=$SCRIPT_PATH"/k8s/"
echo -e "${COLOR_INFO}delete k8s ns${COLOR_NC}"

kubectl delete ns spring-petclinic --force --grace-period=0

echo -e "${COLOR_INFO}delete ingress rule ${COLOR_NC}"
$EGCMD mesh ingress delete ingress1

echo -e "${COLOR_INFO}delete spring petclinic mesh services${COLOR_NC}"
$EGCMD mesh service delete api-gateway
check_result $?
$EGCMD mesh service delete config-server
check_result $?
$EGCMD mesh service delete customers-service
check_result $?
$EGCMD mesh service delete vets-service
check_result $?
$EGCMD mesh service delete visits-service
check_result $?

echo -e "${COLOR_INFO}delete pet tenant ${COLOR_NC}"
$EGCMD mesh tenant delete pet
check_result $?



