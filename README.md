# dnsexit-ip-update
[![CI](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/ci.yml/badge.svg)](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/ci.yml)
[![Release](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/release.yml/badge.svg)](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/release.yml)

A dynamic DNS client for [DNSExit](https://www.dnsexit.com/) registered domains.

This client was built according to the [DNS API Guide](https://dnsexit.com/dns/dns-api/#guide-to-use).  
Before using this client you must create an [DNSExit DNS API key](https://dnsexit.com/dns/dns-api/#apikey).

## Install
#### Binaries
Binaries for Linux, MacOS, and Windows 64-bit architectures can be found on the [releases page](https://github.com/christianTragesser/dnsexit-ip-update/releases).

#### Homebrew Tap
```
brew install christiantragesser/tap/dnsexit
```

#### Container Image
[christiantragesser/dnsexit-ip-update](https://hub.docker.com/r/christiantragesser/dnsexit-ip-update) 

#### Linux Install Package
64-bit architecture DEB and RPM packages can be found on the [releases page](https://github.com/christianTragesser/dnsexit-ip-update/releases).

## Use
```
$ dnsexit -h
Usage of dnsexit:
  -domain string
    	DNSExit domain name
  -interval int
    	Time interval in minutes (default 10)
  -ip string
    	Desired A record IP address
  -key string
    	DNSExit API key
```
#### CLI
```
$ dnsexit -domain <dnsexit domain> -key <API key>
```  
The values for `domain` and `key` can also be configured using environment variables.  
CLI flag values take precedence over environment variable values.
```
$ export DOMAIN="<dnsexit domain>"
$ export API_KEY="<API key>"
$ dnsexit
```
#### Container Instance
```
$ docker run -d christiantragesser/dnsexit-ip-update -domain <dnsexit domain> -key <API key>
``` 
or
```
$ docker run -d -e DOMAIN="<dnsexit domain>" -e API_KEY="<API key>" christiantragesser/dnsexit-ip-update
``` 

### Options
**Check Interval**  
By default, IP update checks happen in 10 minute intervals.  
This cadence can be changed by using the `-interval` flag with a value of the desired interval in minutes.
```
$ dnsexit -domain <dnsexit domain> -key <API key> -interval 20
```  
The `interval` value can also be configured by setting the environment variable `CHECK_INTERVAL`.  

**Preferred IP Address**  
By default, the client configures DNS A record updates using a discovered egress IP address.  
Use the `-ip` flag with a desired IP address to override the discovered IP address value.
```
$ dnsexit -domain <dnsexit domain> -key <API key> -ip 5.5.5.5
```  
The `ip` value can also be configured by setting the environment variable `IP_ADDR`.  