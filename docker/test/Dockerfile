FROM debian:buster-slim
RUN echo "deb http://ftp.cn.debian.org/debian sid main" > /etc/apt/sources.list; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        curl \
        apache2 \
        apache2-utils \
        build-essential \
        libapache2-mod-wsgi-py3 \
        python3 \
        python3-pip \
        python3-dev \
        default-libmysqlclient-dev \
        checkinstall \
        openssl; \
     rm -rf /var/lib/apt/lists/*
RUN  pip install setuptools wheel;\
     apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false;\
     rm -rf /var/lib/apt/lists/*
STOPSIGNAL SIGWINCH
COPY httpd-foreground /usr/local/bin/
CMD ["httpd-foreground"]