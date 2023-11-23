#!/bin/sh
echo "0 6 * * * /usr/local/bin/docker-compose -f /opt/containers/mail-reminder/docker-compose.yml up --build --exit-code-from go-app" > /etc/crontabs/root
crond -f -l 8
