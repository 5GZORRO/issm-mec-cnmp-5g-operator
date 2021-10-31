#!/bin/sh

count=`ps aux | grep "[n]rf" | wc -l`
if [ $count -gt 0 ];then
  pkill -f nrf
fi
