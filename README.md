# dnsexit-ip-update
[![CI](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/ci.yml/badge.svg)](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/ci.yml)
[![Release](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/release.yml/badge.svg)](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/release.yml)

A dynamic DNS client for [DNSExit](https://www.dnsexit.com/) registered domains.

This client was built according to the [DNS API Guide](https://dnsexit.com/dns/dns-api/#guide-to-use).  
Currently, only updates to existing `A` records are supported.  
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
  -domains string
    	DNSExit domain names
  -interval int
    	Time interval in minutes (default 10)
  -ip string
    	Desired A record IP address
  -key string
    	DNSExit API key
```
#### CLI
```
$ dnsexit -domains <dnsexit domain> -key <API key>
```
The values for `domains` and `key` can also be configured using environment variables.  
CLI flag values take precedence over environment variable values.
```
$ export DOMAINS="<dnsexit domain>"
$ export API_KEY="<API key>"
$ dnsexit
```
#### Container Instance
```
$ docker run -d christiantragesser/dnsexit-ip-update -domains <dnsexit domain> -key <API key>
```
or
```
$ docker run -d -e DOMAINS="<dnsexit domain>" -e API_KEY="<API key>" christiantragesser/dnsexit-ip-update
```

### Options
**Multiple DNSExit Domains**  
Multiple DNSExit registered domains can be managed by the same A record information by providing a comma deliniated list of hostnames for the `domains` value.  
```
$ dnsexit -domains my-site.com,your-site.io,our-site.net -key <API key>
```

**Check Interval**  
By default, IP update checks happen in 10 minute intervals.  
This cadence can be changed by using the `-interval` flag with a value of the desired interval in minutes.
```
$ dnsexit -domains <dnsexit domain> -key <API key> -interval 20
```  
The `interval` value can also be configured by setting the environment variable `CHECK_INTERVAL`.  

**Preferred IP Address**  
By default, the client configures DNS A record updates using a discovered egress IP address.  
Use the `-ip` flag with a desired IP address to override the discovered IP address value.
```
$ dnsexit -domain <dnsexit domain> -key <API key> -ip 1.1.1.1
```  
The `ip` value can also be configured by setting the environment variable `IP_ADDR`.  