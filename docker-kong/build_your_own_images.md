# Building Docker images that contain the Kong Gateway

Kong Software uses the [docker-kong github
repository](https://github.com/Kong/docker-kong/) to build Docker images that
contain the Kong Gateway. We will use
[docker-kong](https://github.com/Kong/docker-kong/) as a reference
implementation for describing how to build your own Docker images that contain the Kong
Gateway.

We do not provide a Dockerfile with the `FROM` argument parametrized to allow you to use your
own base image because doing so impedes our ability to get public images
accepted promptly by Dockerhub. If you wish, you can clone the [docker-kong github
repository](https://github.com/Kong/docker-kong/) and adjust the Dockerfile for
your desired package type to use your desired base image and package version,
then use the `build_v2` target in our `Makefile` to build your image. This
document instead takes the approach of walking through the contents of the Dockerfiles
so that you can create and maintain your own.

To build your Docker image, you will need to provide

1. A base image of your choice
1. An entrypoint script that runs the Kong Gateway
1. A Dockerfile that installs the Kong Gateway from a location you specify

## Base image
You can use images derived from RHEL or Ubuntu; Kong Software pushes `.deb`, and
`.rpm` packages to our [public package repository](https://packages.konghq.com/).

## Entrypoint script

Get the [entrypoint
script](https://raw.githubusercontent.com/Kong/docker-kong/master/docker-entrypoint.sh)
from the [docker-kong github repository](https://github.com/Kong/docker-kong/) and put it in
directory where you are planning to run the command to build your Docker image.

## Create a Dockerfile to install Kong Gateway

### Decide how to get the Kong Gateway package
Kong Software provides `.deb`, and `.rpm` packages via our [public package
repository](https://packages.konghq.com/). Decide whether you want your
Dockerfile to

1. Download the desired package from https://packages.konghq.com, or
2. Download the desired package from another package repository you specify, or
3. Install the desired package locally from disk.

If you choose 1 or 2, run the command `touch kong.rpm` in the directory to which
your Dockerfile will download the file; this guarantees that the downloaded file
will have the correct user, groups, and permissions.

If you choose 2 or 3, then download the package you wish to install and put it
in the desired location.

### Write a Dockerfile to install the Kong Gateway package
Use the template below to create your Dockerfile. Angle brackets (`<>`) indicate
values that you need to provide. Comments that start "# Uncomment" indicate that
you need to uncomment lines relevant to your context.

The template is based upon the Dockerfiles in the [docker-kong github
repository](https://github.com/Kong/docker-kong/) and created manually. Check
the Dockerfiles for changes.

```
FROM <your-base-image>

ARG KONG_VERSION=<Kong-Gateway-version>=
ENV KONG_VERSION $KONG_VERSION

# Uncomment the ARG KONG_SHA256 line to build a container using a .deb or .rpm package
# For .deb packages, the SHA is in
# https://cloudsmith.io/~kong/repos/gateway-<gateway-major-version><gateway-minor-version>/packages/detail/deb/kong/<gateway-version>/a=amd64;xc=main;d=debian%252F<os_version>;t=binary/
# For .rpm packages, the SHA is in
# https://cloudsmith.io/~kong/repos/gateway-<gateway-major-version><gateway-minor-version>/packages/detail/rpm/kong/<gateway-version>/a=x86_64;d=el%252F<os_version>;t=binary/
# ARG KONG_SHA256="<.deb-or.rpm-SHA>"

# Uncomment to download package from a remote repository
# ARG ASSET=remote

# Uncomment to install package from local disk
# ARG ASSET=local

ARG EE_PORTS

# Uncomment if you are installing a .rpm
# COPY kong.rpm /tmp/kong.rpm

# Uncomment if you are installing a .deb
# COPY kong.deb /tmp/kong.deb

# Uncomment if you are installing a .deb.tar.gz
# COPY kong.deb.tar.gz /tmp/kong.deb.tar.gz

# hadolint ignore=DL3015
# Uncomment the following section if you are installing a .rpm
# Edit the DOWNLOAD_URL line to install from a repository other than
# packages.konghq.com
# RUN set -ex; \
#     if [ "$ASSET" = "remote" ] ; then \
#       VERSION=$(grep '^VERSION_ID' /etc/os-release | cut -d = -f 2 | sed -e 's/^"//' -e 's/"$//' | cut -d . -f 1) \
#       && KONG_REPO=$(echo ${KONG_VERSION%.*} | sed 's/\.//') \
#       && DOWNLOAD_URL="https://packages.konghq.com/public/gateway-$KONG_REPO/rpm/el/$VERSION/x86_64/kong-$KONG_VERSION.el$VERSION.x86_64.rpm" \
#       && curl -fL $DOWNLOAD_URL -o /tmp/kong.rpm \
#       && echo "$KONG_SHA256  /tmp/kong.rpm" | sha256sum -c -; \
#     fi \
#     && yum install -y /tmp/kong.rpm \
#     && rm /tmp/kong.rpm \
#     && chown kong:0 /usr/local/bin/kong \
#     && chown -R kong:0 /usr/local/kong \
#     && ln -s /usr/local/openresty/bin/resty /usr/local/bin/resty \
#     && ln -s /usr/local/openresty/luajit/bin/luajit /usr/local/bin/luajit \
#     && ln -s /usr/local/openresty/luajit/bin/luajit /usr/local/bin/lua \
#     && ln -s /usr/local/openresty/nginx/sbin/nginx /usr/local/bin/nginx \
#     && kong version

# Uncomment the following section if you are installing a .deb
# Edit the DOWNLOAD_URL line to install from a repository other than
# packages.konghq.com
# RUN set -ex; \
#     apt-get update; \
#     apt-get install -y curl; \
#     if [ "$ASSET" = "remote" ] ; then \
#       CODENAME=$(cat /etc/os-release | grep VERSION_CODENAME | cut -d = -f 2) \
#       && KONG_REPO=$(echo ${KONG_VERSION%.*} | sed 's/\.//') \
#       && DOWNLOAD_URL="https://packages.konghq.com/public/gateway-$KONG_REPO/deb/ubuntu/pool/$CODENAME/main/k/ko/kong_$KONG_VERSION/kong_${KONG_VERSION}_amd64.deb" \
#       && curl -fL $DOWNLOAD_URL -o /tmp/kong.deb \
#       && echo "$KONG_SHA256  /tmp/kong.deb" | sha256sum -c -; \
#     fi \
#     && apt-get update \
#     && apt-get install --yes /tmp/kong.deb \
#     && rm -rf /var/lib/apt/lists/* \
#     && rm -rf /tmp/kong.deb \
#     && chown kong:0 /usr/local/bin/kong \
#     && chown -R kong:0 /usr/local/kong \
#     && ln -s /usr/local/openresty/bin/resty /usr/local/bin/resty \
#     && ln -s /usr/local/openresty/luajit/bin/luajit /usr/local/bin/luajit \
#     && ln -s /usr/local/openresty/luajit/bin/luajit /usr/local/bin/lua \
#     && ln -s /usr/local/openresty/nginx/sbin/nginx /usr/local/bin/nginx \
#     && kong version \
#     && apt-get purge curl -y

COPY docker-entrypoint.sh /docker-entrypoint.sh

USER kong

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 8000 8443 8001 8444 $EE_PORTS

STOPSIGNAL SIGQUIT

HEALTHCHECK --interval=60s --timeout=10s --retries=10 CMD kong-health

CMD ["kong", "docker-start"]
```

### Run the docker command

Run the command `docker build --no-cache -t kong-<your-base-image>
<path-for-built-image>` to build the docker image.
