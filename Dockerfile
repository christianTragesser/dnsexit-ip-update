FROM python:3-alpine

COPY requirements/requirements.txt /tmp/requirements.txt

RUN apk add ca-certificates --no-cache && \
    pip install -r /tmp/requirements.txt --no-cache-dir && \
    addgroup -S -g 2222 dnsexit && \
    adduser -S -u 2222 -g dnsexit dnsexit

COPY src/*.py /opt/

RUN chmod 755 -R /opt

USER dnsexit

CMD ["/bin/sh", "-c", "python /opt/main.py"]