#!/bin/bash

usage() {
   echo "$0 features:" >&2
   echo '  - If set, use $SSODOMAIN environment variable to set domain name' >&2
   echo '  - setup llng-fastcgi-server envs and dirs' >&2
}

if [ "$1" = '--help' ]; then
    usage
    exit 1
fi

if [ -n "$SSODOMAIN" ]; then
    echo "Using ServerName: $SSODOMAIN"

    sed -i "s/example\.com/${SSODOMAIN}/g" \
        /etc/lemonldap-ng/* \
        /var/lib/lemonldap-ng/conf/lmConf-1.js* \
        /var/lib/lemonldap-ng/test/index.pl
fi

# Setup like lemonldap-ng-fastcgi-server.service
. /etc/default/lemonldap-ng-fastcgi-server
export SOCKET
export PID
export USER
export GROUP
mkdir -p "$(dirname "$PID")"
chown "$USER" "$(dirname "$PID")"

exec "$@"
