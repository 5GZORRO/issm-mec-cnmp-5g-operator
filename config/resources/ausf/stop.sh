#!/bin/sh

count=`ps aux | grep "[a]usf" | wc -l`
if [ $count -gt 0 ];then
  pkill -f ausf
fi
