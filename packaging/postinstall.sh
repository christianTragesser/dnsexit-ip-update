#!/bin/sh
set -e
rm -f /usr/local/bin/dnsexit
ln -s /usr/local/bin/dnsexit-linux-* /usr/local/bin/dnsexit
