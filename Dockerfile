FROM alpine:latest

RUN apk add curl perl ca-certificates --no-cache && \
  curl http://downloads.dnsexit.com/ipUpdate-1.71.tar.gz | tar -zx -C /opt && \
  cp /opt/dnsexit/Http_get.pm /usr/lib/perl5/core_perl/

ENV LOGIN="" PASS="" DOMAIN=""

COPY dnsexit.conf /etc/
COPY ipUpdate.sh /usr/local/bin/

ENTRYPOINT ["ipUpdate.sh"]
CMD ["sh"]