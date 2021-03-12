#!/usr/bin/env bash
cd /var/www/html/short/example
git pull ; go build ; systemctl restart short ; systemctl status short