#!/bin/sh

set -e

service="$1"
port="$2"
shift

echo Waiting for kafka service start...;
while ! nc -z $service $port;
do
  sleep 1;
done;
echo Connected!;
./compiled-app
