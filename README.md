# dnsexit-ip-update
[![pipeline status](https://gitlab.com/christianTragesser/dnsexit-ip-update/badges/master/pipeline.svg)](https://gitlab.com/christianTragesser/dnsexit-ip-update/commits/master)

A python dynamic DNS client for [DNSExit](https://www.dnsexit.com/) registered domains.

Before using this client it is **strongly recommended** you create a [Dynamic IP Update Password](https://www.dnsexit.com/Direct.sv?cmd=userProfilePwIP) for your account rather than using your DNSExit account login credentials.

This client was built according to the DNSExit IP Update [specification document](http://downloads.dnsexit.com/ipUpdateDev.doc).

#### Install
##### [PyPi](https://pypi.org/project/dnsexit-ip-update/)
Python 3.6 or later
```sh
$ export LOGIN="<your dnsexit login>"
$ export PASSWORD="<your dnsexit IP Update password>"
$ export DOMAIN="<your dnsexit registered domain>"
$ pip install dnsexit-ip-update
$ python -m dnsexitUpdate
```
##### Docker (suggested)
```sh
$ docker run -d -e LOGIN="<your dnsexit login>" \
                -e PASSWORD="<your dnsexit IP Update password>" \
                -e DOMAIN="<your dnsexit registered domain>" \
                registry.gitlab.com/christiantragesser/dnsexit-ip-update
```
