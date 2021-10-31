#!/bin/sh

count=`ps aux | grep "[n]ssf" | wc -l`
if [ $count -gt 0 ];then
  pkill -f nssf
fi
