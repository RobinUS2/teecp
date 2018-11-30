FROM ubuntu:18.04
MAINTAINER Loc Tran <loc@route42.nl>

# Update
RUN apt-get update && \
    apt-get install -y golang && \
    apt-get dist-upgrade -y && \
    apt-get install -y ca-certificates && \
    apt-get install -y gcc-multilib && \
    apt-get install -y gcc-mingw-w64 && \
    apt-get install -y git && \
    apt-get install -y libpcap-dev

# Config
VOLUME ["/usr/local/src"]

RUN mkdir -p /usr/local/app
ENV GOPATH=/usr/local/app/gopath/
ENV GOBIN=/usr/local/app/src/

WORKDIR /usr/local/app/src/
