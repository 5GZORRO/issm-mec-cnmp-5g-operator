#!/bin/sh

count=`ps aux | grep "[m]ongo" | wc -l`
if [ $count -eq 0 ];then
  echo "Starting Mongo..."
  mongod --port 27017 --bind_ip 0.0.0.0 > /dev/null 2> /dev/null &
  echo "...done"
fi
