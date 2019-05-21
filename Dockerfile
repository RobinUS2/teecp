FROM ubuntu:18.04
MAINTAINER Loc Tran <loc@route42.nl>

# Update
RUN apt-get update && \
    apt-get dist-upgrade -y

RUN apt-get update && \
    apt-get install software-properties-common -y && \
    add-apt-repository ppa:longsleep/golang-backports && \
    apt-get update && \
    apt-get install -y golang-1.12 && \
    apt-get install -y ca-certificates && \
    apt-get install -y gcc-multilib && \
    apt-get install -y gcc-mingw-w64 && \
    apt-get install -y git && \
    apt-get install -y libpcap-dev

# Config
VOLUME ["/usr/local/src"]

RUN mkdir -p /usr/local/app
ENV GOPATH=/root
ENV GOBIN=/usr/local/app/forwarder/

WORKDIR /usr/local/app/
