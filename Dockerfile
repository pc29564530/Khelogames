FROM python:3.11-slim-bullseye

MAINTAINER Jookies LTD <jasmin@jookies.net>

# add our user and group first to make sure their IDs get assigned consistently, regardless of whatever dependencies get added
RUN groupadd -r jasmin && useradd -r -g jasmin jasmin

# Install requirements
RUN apt-get update && apt-get install -y \
    libffi-dev \
    libssl-dev \
    # Run python with jemalloc
    # More on this:
    # - https://zapier.com/engineering/celery-python-jemalloc/
    # - https://paste.pics/581cc286226407ab0be400b94951a7d9
    libjemalloc2

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Run python with jemalloc
ENV LD_PRELOAD /usr/lib/aarch64-linux-gnu/libjemalloc.so.2

# Install Jasmin SMS gateway
RUN mkdir -p /etc/jasmin/resource \
    /etc/jasmin/store \
    /var/log/jasmin \
  && chown jasmin:jasmin /etc/jasmin/store \
    /var/log/jasmin

WORKDIR /build

COPY . .

RUN pip install .

ENV UNICODEMAP_JP unicode-ascii

ENV ROOT_PATH /
ENV CONFIG_PATH /etc/jasmin
ENV RESOURCE_PATH /etc/jasmin/resource
ENV STORE_PATH /etc/jasmin/store
ENV LOG_PATH /var/log/jasmin

COPY misc/config/*.cfg ${CONFIG_PATH}/
COPY misc/config/resource ${RESOURCE_PATH}

WORKDIR /usr/jasmin

# Change binding host for jcli, redis, and amqp
RUN sed -i '/\[jcli\]/a bind=0.0.0.0' /etc/jasmin/jasmin.cfg
RUN sed -i '/\[redis-client\]/a host=redis' /etc/jasmin/jasmin.cfg
RUN sed -i '/\[amqp-broker\]/a host=rabbitmq' /etc/jasmin/jasmin.cfg

EXPOSE 2775 8990 1401
VOLUME ["/var/log/jasmin", "/etc/jasmin", "/etc/jasmin/store"]

COPY docker/docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["jasmind.py", "--enable-interceptor-client", "--enable-dlr-thrower", "--enable-dlr-lookup", "-u", "jcliadmin", "-p", "jclipwd"]
# Notes:
# - jasmind is started with native dlr-thrower and dlr-lookup threads instead of standalone processes
# - restapi (0.9rc16+) is not started in this docker configuration
