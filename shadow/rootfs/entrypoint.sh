#!/bin/sh

if [ "$(echo $1 | head -c 1)" != "-" ] ; then
  exec "$@"
else
  exec /opt/easemesh-shadowservice-controller/bin/easemesh-shadowservice-controller "$@"
fi

