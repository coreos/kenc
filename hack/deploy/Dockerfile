FROM debian:jessie

RUN DEBIAN_FRONTEND=noninteractive apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get -yy -q install \
    iptables \
    ca-certificates \
    && DEBIAN_FRONTEND=noninteractive apt-get upgrade -y \
    && DEBIAN_FRONTEND=noninteractive apt-get autoremove -y \
    && DEBIAN_FRONTEND=noninteractive apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* # CACHEBUST

ADD _output/bin/kenc /usr/local/bin

CMD ["/usr/local/bin/kenc"]
