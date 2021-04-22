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

echo -e "${COLOR_INFO}create pet tenant ${COLOR_NC}"
$EGCMD mesh tenant create  -f $CONFIG_PATH"/pet-tenant.yaml"
check_result $?

echo -e "${COLOR_INFO}create spring petclinic mesh services${COLOR_NC}"
$EGCMD mesh service create -f $CONFIG_PATH"/api-gateway.yaml"
check_result $?
$EGCMD mesh service create -f $CONFIG_PATH"/config-server.yaml"
check_result $?
$EGCMD mesh service create -f $CONFIG_PATH"/customers-service.yaml"
check_result $?
$EGCMD mesh service create -f $CONFIG_PATH"/vets-service.yaml"
check_result $?
$EGCMD mesh service create -f $CONFIG_PATH"/visits-service.yaml"
check_result $?

echo -e "${COLOR_INFO}create ingress rule ${COLOR_NC}"
$EGCMD mesh ingress create -f $CONFIG_PATH"/ingress-rule.yaml"


K8S_YAML_PATH=$SCRIPT_PATH"/k8s/"
echo -e "${COLOR_INFO}create k8s ns${COLOR_NC}"
kubectl create ns spring-petclinic  

echo -e "${COLOR_INFO}deploy k8s mesh deployment ${COLOR_NC}"
kubectl apply -f $K8S_YAML_PATH
