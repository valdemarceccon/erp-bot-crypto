#!/bin/bash

printenv | grep POSTGRES > /etc/environment
crontab /etc/cron.d/crontab
cron -f
