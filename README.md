# dnsexit-ip-update
[![CI](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/ci.yml/badge.svg)](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/ci.yml)
[![Release](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/release.yml/badge.svg)](https://github.com/christianTragesser/dnsexit-ip-update/actions/workflows/release.yml)

A dynamic DNS client for [DNSExit](https://www.dnsexit.com/) registered domains.

This client was built according to the [DNS API Guide](https://dnsexit.com/dns/dns-api/#guide-to-use).  
Before using this client you must create a [DNSExit DNS API key](https://dnsexit.com/dns/dns-api/#apikey).

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
        DNSExit domain name(s), comma-separated
  -key string
        DNSExit API key
  -ip string
        Desired record IP address (auto-discovered if omitted)
  -record-type string
        DNS record type: A, AAAA, or SELF (default: A)
  -ttl int
        Record TTL in minutes (default 5)
  -interval int
        Update check interval in minutes, minimum 5 (default 10)
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
Multiple DNSExit registered domains can be managed with the same record settings by providing a comma-separated list of hostnames for the `domains` value.
```
$ dnsexit -domains my-site.com,your-site.io,our-site.net -key <API key>
```

**Record Type**  
Three record types are supported. The default is `A` (IPv4).

| Type | Description |
|------|-------------|
| `A` | IPv4 address record (default) |
| `AAAA` | IPv6 address record |
| `SELF` | DNSExit detects the client's public IP server-side; no IP address required |

```
$ dnsexit -domains <dnsexit domain> -key <API key> -record-type AAAA -ip 2001:db8::1
```

The `SELF` type is the simplest option for IPv4 dynamic DNS — DNSExit reads the client's public IP from the incoming request, so no IP discovery or `-ip` flag is needed:
```
$ dnsexit -domains <dnsexit domain> -key <API key> -record-type SELF
```

The `record-type` value can also be configured by setting the environment variable `RECORD_TYPE`.

**Record TTL**  
The record TTL controls how long resolvers cache the DNS record, in minutes. The default is `5` minutes, which allows IP changes to propagate quickly. Use a higher value if your IP address is stable and you want to reduce DNS query load.

```
$ dnsexit -domains <dnsexit domain> -key <API key> -ttl 30
```

A TTL of `0` is valid and is used by some certificate issuance workflows (e.g. Let's Encrypt DNS challenges).

The `ttl` value can also be configured by setting the environment variable `RECORD_TTL`.

**Check Interval**  
By default, IP update checks happen in 10 minute intervals. The DNSExit API requires a minimum interval of 5 minutes; values below this are clamped automatically.

```
$ dnsexit -domains <dnsexit domain> -key <API key> -interval 20
```

The `interval` value can also be configured by setting the environment variable `CHECK_INTERVAL`.

**Preferred IP Address**  
By default, the client discovers the egress IP address automatically and re-checks it on every interval. Use the `-ip` flag to pin a specific IP address instead. The provided value is validated as a correct IP address before the client starts.

```
$ dnsexit -domains <dnsexit domain> -key <API key> -ip 1.1.1.1
```

The `ip` value can also be configured by setting the environment variable `IP_ADDR`.

**Subdomains**  
Subdomain records are supported. Provide the full hostname — the client correctly splits the base domain and subdomain when communicating with the DNSExit API.

```
$ dnsexit -domains api.example.com -key <API key>
```

### Environment Variable Reference

| Variable | Flag equivalent | Description |
|----------|----------------|-------------|
| `DOMAINS` | `-domains` | Comma-separated list of domain names |
| `API_KEY` | `-key` | DNSExit API key |
| `RECORD_TYPE` | `-record-type` | DNS record type: `A`, `AAAA`, or `SELF` |
| `IP_ADDR` | `-ip` | Pinned IP address (auto-discovered if unset) |
| `RECORD_TTL` | `-ttl` | Record TTL in minutes (default `5`) |
| `CHECK_INTERVAL` | `-interval` | Update check interval in minutes (default `10`, minimum `5`) |
