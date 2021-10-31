#!/bin/sh

count=`ps aux | grep "[p]cf" | wc -l`
if [ $count -gt 0 ];then
  pkill -f pcf
fi
