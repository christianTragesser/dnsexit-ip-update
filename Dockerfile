FROM python:3-alpine

COPY requirements/requirements.txt /tmp/requirements.txt

RUN apk add ca-certificates --no-cache && \
    pip install -r /tmp/requirements.txt --no-cache-dir && \
    addgroup -S -g 2222 dnsexit && \
    adduser -S -u 2222 -g dnsexit dnsexit && \
    mkdir /opt/dnsexitUpdate

COPY src/dnsexitUpdate/*.py /opt/dnsexitUpdate/

RUN chmod 755 -R /opt/dnsexitUpdate

USER dnsexit

WORKDIR /opt

CMD ["/bin/sh", "-c", "python -m dnsexitUpdate"]