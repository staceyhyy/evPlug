#!/usr/bin/env bash

if [ -n "$1" ]; then
    if [ $1 = 'rebuild' ]; then 
        docker-compose down -v
        docker-compose build --no-cache
        docker-compose -f ./docker-compose.yml up -d --remove-orphans
    elif [ $1 = 'restart' ]; then 
        docker-compose restart
    fi
else
    echo "Usage: setup.sh [rebuild|restart]"
    echo 
    echo "rebuild   tear down and rebuild all containers"
    echo "restart   stop and restart all containers"
fi