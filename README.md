# dnsexit-ip-update
[![pipeline status](https://gitlab.com/christianTragesser/dnsexit-ip-update/badges/master/pipeline.svg)](https://gitlab.com/christianTragesser/dnsexit-ip-update/commits/master)

A dynamic DNS client for [DNSExit](https://www.dnsexit.com/) registered domains.

Before using this client it is **strongly recommended** you create a [Dynamic IP Update Password](https://www.dnsexit.com/Direct.sv?cmd=userProfilePwIP) for your account rather than using your DNSExit account login credentials.

This client was built according to the DNSExit IP Update [specification document](http://downloads.dnsexit.com/ipUpdateDev.doc).

### Install
#### [PyPi](https://pypi.org/project/dnsexit-ip-update/)
For systems using Python 3.6 or later, there is pip package available:
```sh
$ pip install dnsexit-ip-update
```
#### [Docker](https://gitlab.com/christianTragesser/dnsexit-ip-update/container_registry) (suggested)
This package is available as a docker image as well.

`registry.gitlab.com/christiantragesser/dnsexit-ip-update`

### Configure and Run
#### Python Package
```sh
$ export LOGIN="<your dnsexit login>"
$ export PASSWORD="<your dnsexit IP Update password>"
$ export DOMAIN="<your dnsexit registered domain>"
$ python -m dnsexitUpdate
```
#### Docker
```sh
$ docker run -d -e LOGIN="<your dnsexit login>" \
                -e PASSWORD="<your dnsexit IP Update password>" \
                -e DOMAIN="<your dnsexit registered domain>" \
                registry.gitlab.com/christiantragesser/dnsexit-ip-update
```

### Configure Options
**Check Interval**  
By default IP update checks happen in 10 minute intervals.  This cadence can be changed by setting the enviromental variable `CHECK_INTERVAL` to the desired interval in units of seconds.
```sh
# 20 minute interval
export CHECK_INTERVAL=1200
```