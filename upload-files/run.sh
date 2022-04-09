#!/bin/bash
##########################################################
# File Name: build.sh
# Author: ims
# Created Time: 2020年10月06日 星期二 11时00分33秒
##########################################################

#ubuntu root lack mkdocs resolve
# vi /usr/local/bin/mkdocs
# add follow context
# sys.path.append('/home/username/.local/lib/python3.8/site-packages')
#daemonize 原文链接：https://blog.csdn.net/erlang_hell/article/details/51187205

PID_FILE=daemonize.pid
PORT=9000
ERR_LOG=error.log
STD_LOG=stdout.log

if [ $# -eq 0 ]; then
    go build main.go
    echo "daemonize run"
    echo "" > $ERR_LOG
    echo "" > $STD_LOG
    daemonize -a -e ./$ERR_LOG -o $STD_LOG -p $PID_FILE -l $PID_FILE -c ./  `pwd`/main
    tail -f $STD_LOG
    exit
fi

if [[ $1 =~ "g" ]];then
    echo "listen on 80 port, global daemon"
    nohup mkdocs serve -a 0.0.0.0:80 &
fi

if [[ $1 =~ "local" ]];then
    echo "local mode"
    mkdocs serve
fi

if [[ $1 =~ "kill" ]];then
    pid=`cat $PID_FILE`
    echo pid:$pid
    kill $pid
fi
