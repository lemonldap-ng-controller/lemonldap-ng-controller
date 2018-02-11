#!/bin/sh

# Copy static files if needed
if [ -d /usr/share/lemonldap-ng/portal-skins \
  -a -d /srv/var/lib/lemonldap-ng/portal/skins \
  -a -z "$(ls /srv/var/lib/lemonldap-ng/portal/skins)" \
]; then
  cp -aT /usr/share/lemonldap-ng/portal-skins /srv/var/lib/lemonldap-ng/portal/skins
fi

# Setup like lemonldap-ng-fastcgi-server.service
. /etc/default/lemonldap-ng-fastcgi-server
export SOCKET PID USER GROUP
mkdir -p "$(dirname "$PID")"
chown "$USER" "$(dirname "$PID")"

exec "$@"
