#!/bin/sh

count=`ps aux | grep "[s]mf" | wc -l`
if [ $count -gt 0 ];then
  pkill -f smf
fi
