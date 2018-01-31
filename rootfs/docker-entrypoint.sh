#!/bin/sh

# Setup like lemonldap-ng-fastcgi-server.service
. /etc/default/lemonldap-ng-fastcgi-server
export SOCKET PID USER GROUP
mkdir -p "$(dirname "$PID")"
chown "$USER" "$(dirname "$PID")"

exec "$@"
