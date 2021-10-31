#!/bin/sh

count=`ps aux | grep "[a]mf" | wc -l`
if [ $count -gt 0 ];then
  pkill -f amf
fi
