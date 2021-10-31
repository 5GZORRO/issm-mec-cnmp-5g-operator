#!/bin/sh

count=`ps aux | grep "[m]ongo" | wc -l`
if [ $count -gt 0 ];then
  pkill mongod
fi
