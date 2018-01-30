FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

COPY docker-entrypoint.sh lemonldap-ng-controller lemonldap-ng.key lemonldap-ng.list /

RUN set -x && \
    echo "# Adding deb repository..." && \
    apt-get update && \
    apt-get dist-upgrade -y && \
    apt-get install -y \
      apt-transport-https \
      gnupg \
      && \
    mv /lemonldap-ng.list /etc/apt/sources.list.d/ && \
    apt-key add /lemonldap-ng.key && \
    rm /lemonldap-ng.key && \
    apt-get update && \
    \
    echo "# Installing packages..." && \
    apt-get install -y \
      dumb-init \
      lemonldap-ng-fastcgi-server \
      && \
    \
    echo "# Cleaning up..." && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 9000

ENTRYPOINT ["dumb-init","--","/docker-entrypoint.sh"]

CMD "/lemonldap-ng-controller" "--alsologtostderr"
