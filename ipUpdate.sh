#!/bin/sh

i=0

if [ -z "$LOGIN" ] || [ -z "$PASS" ] || [ -z "$DOMAIN" ]; then
  echo -e "\n    Missing required login, password, or domain env var - exiting...\n"
  exit 1;
fi

sed -i -e "s/LOGIN_HERE/$LOGIN/; s/PASS_HERE/$PASS/; s/DOMAIN_HERE/$DOMAIN/" /etc/dnsexit.conf

while [ $i -eq 0 ]; do
  /opt/dnsexit/ipUpdate.pl
  sleep 600
done
